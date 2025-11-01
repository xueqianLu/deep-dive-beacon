package services

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BeaconBlockService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
}

func NewBeaconBlockService(db *gorm.DB, redis *redis.Client, logger *logrus.Logger) *BeaconBlockService {
	return &BeaconBlockService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}
