package beaconscanner

import (
	"fmt"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	lru "github.com/hashicorp/golang-lru"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	beaconapi "github.com/xueqianLu/deep-dive-beacon/beacon"
	"github.com/xueqianLu/deep-dive-beacon/config"
	"github.com/xueqianLu/deep-dive-beacon/constant"
	"github.com/xueqianLu/deep-dive-beacon/internal/services"
	"github.com/xueqianLu/deep-dive-beacon/models/dbmodels"
	"gorm.io/gorm"
	"math/big"
	"sync"
	"time"
)

var (
	big0 = big.NewInt(0)
)

type DirectlyBlockScanner struct {
	config       *config.Config
	db           *gorm.DB
	rdb          *redis.Client
	logger       *logrus.Logger
	services     *services.Services
	rwmux        sync.RWMutex
	start        int64
	end          int64
	quit         chan struct{}
	cache        *lru.Cache
	beaconClient *beaconapi.BeaconClient
	running      map[uint]bool
}

func NewDirectlyBlockScanner(cfg *config.Config, db *gorm.DB, redis *redis.Client, logger *logrus.Logger, start int64, end int64) *DirectlyBlockScanner {
	svc := services.NewServices(db, redis, logger, cfg)
	cache, err := lru.New(1000)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create LRU cache")
	}

	scan := &DirectlyBlockScanner{
		config:       cfg,
		db:           db,
		rdb:          redis,
		logger:       logger,
		services:     svc,
		start:        start,
		end:          end,
		quit:         make(chan struct{}),
		running:      make(map[uint]bool),
		cache:        cache,
		beaconClient: beaconapi.NewBeaconGwClient(cfg.Chain.BeaconURL),
	}
	return scan
}

func (s *DirectlyBlockScanner) checkTaskRunning(task *dbmodels.DirectlyScanTask) bool {
	s.rwmux.RLock()
	defer s.rwmux.RUnlock()
	if running, exist := s.running[task.ID]; exist && running {
		return true
	}
	return false
}

func (s *DirectlyBlockScanner) setTaskRunning(task *dbmodels.DirectlyScanTask, running bool) {
	s.rwmux.Lock()
	defer s.rwmux.Unlock()
	s.running[task.ID] = running
}

func (s *DirectlyBlockScanner) Start() error {
	s.logger.Info("Starting blockchain scanner service")

	ticker := time.NewTicker(10 * time.Second) // Scan every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-s.quit:
			s.logger.Info("Scanner service stopped")
			return nil

		case <-ticker.C:
			tasks, err := s.services.DirectlyScan.GetScanTaskByType(constant.DIRECTLY_SCAN_TYPE_BEACON_BLOCK)
			if len(tasks) == 0 || err != nil {
				s.logger.WithError(err).Info("Beacon block scan task is not enabled, skipping...")
				continue
			}
			for _, task := range tasks {
				if s.checkTaskRunning(task) {
					continue
				}
				go func(t *dbmodels.DirectlyScanTask) {
					s.setTaskRunning(t, true)
					defer s.setTaskRunning(t, false)
					if err := s.doScanTask(t); err != nil {
						s.logger.WithError(err).Error("Directly block scan task failed")
					}
				}(task)
			}
			ticker.Reset(time.Second * 10)
		}
	}
}

func (s *DirectlyBlockScanner) Stop() {
	close(s.quit)
}
func intToStr(num uint64) string {
	return fmt.Sprintf("%d", num)
}

func (s *DirectlyBlockScanner) SetFailed(slot int64) {
	key := fmt.Sprintf("failed_%d", slot)
	if _, exist := s.cache.Get(key); exist {
		return
	} else {
		s.cache.Add(key, time.Now().Unix())
	}
}

func (s *DirectlyBlockScanner) ShouldSkip(slot int64) bool {
	key := fmt.Sprintf("failed_%d", slot)
	if val, exist := s.cache.Get(key); exist {
		tm := val.(int64)
		// skip block if failed unless 60 seconds have passed
		if time.Now().Unix()-tm > 10 {
			return true
		}
	}
	return false

}

func (s *DirectlyBlockScanner) doScanTask(task *dbmodels.DirectlyScanTask) error {
	logger := s.logger.WithField("task", task.ID)
	height := task.LastNumber + 1
	if task.LastNumber < task.Start {
		height = task.Start
	}

	if height > task.End {
		logger.Info("Scan task already completed")
		return nil
	}

	var latest *apiv1.BeaconBlockHeader
	var err error

	var refreshLatest = func() {
		latest, err = s.beaconClient.GetLatestBeaconHeader()
		if err != nil {
			logger.WithError(err).Error("Failed to get latest beacon header")
			return
		}
	}

	var running bool = true
	for running {
		if latest == nil {
			refreshLatest()
			continue
		}
		if height > uint64(latest.Header.Message.Slot) {
			time.Sleep(time.Second)
			refreshLatest()
			continue
		}

		block, err := s.beaconClient.GetBlockById(intToStr(height))
		if err != nil {
			logger.WithFields(logrus.Fields{
				"height": height,
			}).WithError(err).Error("Failed to get beacon block by id")
			if s.ShouldSkip(int64(height)) {
				logger.WithField("height", height).Warning("Skipping failed height")
				height++
			} else {
				s.SetFailed(int64(height))
			}
			continue
		}
		tx := s.db.Begin()
		if tx.Error != nil {
			logger.WithError(tx.Error).Error("Failed to begin db transaction")
			return tx.Error
		}
		if err = s.processBeaconBlock(tx, block); err != nil {
			tx.Rollback()
			logger.WithError(err).Error("processor beacon block failed")
			return err
		} else {
			if err := tx.Commit().Error; err != nil {
				logger.WithError(err).Error("commit transaction to db failed")
				return err
			}
		}
		// update task last processed height
		task.LastNumber = height
		s.services.DirectlyScan.UpdateScanTask(task)
		if height%100 == 0 {
			logger.WithFields(logrus.Fields{
				"remain": task.End - task.LastNumber,
				"height": height,
			}).Info("Processed beacon blocks")
		}
		height++
	}
	return nil
}

func (s *DirectlyBlockScanner) processBeaconBlock(db *gorm.DB, blk *spec.VersionedSignedBeaconBlock) error {
	dbblk, err := s.ToDBBlock(blk)
	if err != nil {
		return err
	}
	db.Model(&dbmodels.BeaconBlock{}).Save(dbblk)
	atts := s.GetBlkAtts(blk)
	for _, att := range atts {
		db.Model(&dbmodels.BeaconAttestation{}).Save(att)
	}
	return nil
}

var (
	slotsPerEpoch = uint64(32)
)

func (s *DirectlyBlockScanner) ToDBBlock(blk *spec.VersionedSignedBeaconBlock) (*dbmodels.BeaconBlock, error) {
	slot, _ := blk.Slot()
	dbBlk := new(dbmodels.BeaconBlock)
	dbBlk.SlotNumber = uint64(slot)
	dbBlk.EpochNumber = uint64(slot) / slotsPerEpoch
	if blk.Phase0 != nil {
		// phase0ToDBBlock
		return s.phase0ToDBBlock(blk.Phase0, dbBlk)
	}
	if blk.Altair != nil {
		return s.altairToDBBlock(blk.Altair, dbBlk)
	}
	if blk.Bellatrix != nil {
		return s.bellatrixToDBBlock(blk.Bellatrix, dbBlk)
	}
	if blk.Capella != nil {
		return s.capellaToDBBlock(blk.Capella, dbBlk)
	}
	if blk.Deneb != nil {
		return s.denebToDBBlock(blk.Deneb, dbBlk)
	}
	if blk.Electra != nil {
		return s.electraToDBBlock(blk.Electra, dbBlk)
	}
	if blk.Fulu != nil {
		return s.fuluToDBBlock(blk.Fulu, dbBlk)
	}
	return nil, fmt.Errorf("unknown block version at slot %d", slot)
}

func (s *DirectlyBlockScanner) GetBlkAtts(blk *spec.VersionedSignedBeaconBlock) []*dbmodels.BeaconAttestation {
	if blk.Phase0 != nil {
		return s.getPhase0Attestations(blk.Phase0)
	}
	if blk.Altair != nil {
		return s.getAltairAttestations(blk.Altair)
	}
	if blk.Bellatrix != nil {
		return s.getBellatrixAttestations(blk.Bellatrix)
	}
	if blk.Capella != nil {
		return s.getCapellaAttestations(blk.Capella)
	}
	if blk.Deneb != nil {
		return s.getDenebAttestations(blk.Deneb)
	}
	if blk.Electra != nil {
		return s.getElectraAttestations(blk.Electra)
	}
	if blk.Fulu != nil {
		return s.getFuluAttestations(blk.Fulu)
	}
	return []*dbmodels.BeaconAttestation{}
}
