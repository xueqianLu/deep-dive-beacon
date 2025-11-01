package beaconscanner

import (
	"context"
	"fmt"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
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

type BeaconBlockScanner struct {
	config       *config.Config
	db           *gorm.DB
	rdb          *redis.Client
	logger       *logrus.Logger
	services     *services.Services
	rwmux        sync.RWMutex
	quit         chan struct{}
	cache        sync.Map // string -> TaskContent
	beaconClient *beaconapi.BeaconClient
	running      bool
}

func NewBeaconBlockScanner(cfg *config.Config, db *gorm.DB, redis *redis.Client, logger *logrus.Logger) *BeaconBlockScanner {
	svc := services.NewServices(db, redis, logger, cfg)

	scan := &BeaconBlockScanner{
		config:       cfg,
		db:           db,
		rdb:          redis,
		logger:       logger,
		services:     svc,
		quit:         make(chan struct{}),
		running:      false,
		beaconClient: beaconapi.NewBeaconGwClient(cfg.Chain.BeaconURL),
	}
	return scan
}

func (s *BeaconBlockScanner) Start() error {
	s.logger.Info("Starting blockchain scanner service")

	ticker := time.NewTicker(10 * time.Second) // Scan every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-s.quit:
			s.logger.Info("Scanner service stopped")
			return nil

		case <-ticker.C:

			task, err := s.services.ScanTask.GetScanTaskByType(constant.SCAN_TYPE_BEACON_BLOCK)
			if task == nil || err != nil {
				s.logger.WithError(err).Info("Beacon block scan task is not enabled, skipping...")
				continue
			}
			if !s.running {
				go s.doScanTask(task)
			}
			ticker.Reset(time.Minute)
		}
	}
}

func (s *BeaconBlockScanner) Stop() {
	close(s.quit)
}
func intToStr(num uint64) string {
	return fmt.Sprintf("%d", num)
}

func (s *BeaconBlockScanner) doScanTask(task *dbmodels.ScanTask) error {
	logger := s.logger.WithField("module", "block-scanner")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	height := task.LastNumber + 1

	s.running = true
	defer func() {
		s.running = false
	}()
	tm := time.NewTicker(time.Millisecond * 10)
	defer tm.Stop()

	var latest *apiv1.BeaconBlockHeader
	var err error

	var refreshLatest = func() {
		latest, err = s.beaconClient.GetLatestBeaconHeader()
		if err != nil {
			logger.WithError(err).Error("Failed to get latest beacon header")
			return
		}
	}

	for {
		select {
		case <-ctx.Done():

		case <-s.quit:
			return nil
		case <-tm.C:
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
				}).WithError(err).Error("Failed to get beacon block header by id")
				continue
			}
			tx := s.db.Begin()
			if tx.Error != nil {
				logger.WithError(tx.Error).Error("Failed to begin db transaction")
				continue
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
			s.services.ScanTask.UpdateScanTask(task)
			height++
		}
	}
}

func (s *BeaconBlockScanner) processBeaconBlock(db *gorm.DB, blk *spec.VersionedSignedBeaconBlock) error {
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

func (s *BeaconBlockScanner) ToDBBlock(blk *spec.VersionedSignedBeaconBlock) (*dbmodels.BeaconBlock, error) {
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

func (s *BeaconBlockScanner) GetBlkAtts(blk *spec.VersionedSignedBeaconBlock) []*dbmodels.BeaconAttestation {
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
