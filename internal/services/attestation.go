package services

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttestationService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
}

func NewAttestationService(db *gorm.DB, redis *redis.Client, logger *logrus.Logger) *AttestationService {
	return &AttestationService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}
