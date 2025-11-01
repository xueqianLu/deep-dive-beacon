package services

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/deep-dive-beacon/config"
	"gorm.io/gorm"
)

type Services struct {
	BeaconBlock *BeaconBlockService
	Attest      *AttestationService
	ScanTask    *ScanTaskService
}

func NewServices(db *gorm.DB, redis *redis.Client, logger *logrus.Logger, cfg *config.Config) *Services {
	return &Services{
		BeaconBlock: NewBeaconBlockService(db, redis, logger),
		Attest:      NewAttestationService(db, redis, logger),
		ScanTask:    NewScanTaskService(db, redis, logger),
	}
}
