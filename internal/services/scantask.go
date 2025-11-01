package services

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/deep-dive-beacon/models/dbmodels"
	"gorm.io/gorm"
)

type ScanTaskService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
}

func NewScanTaskService(db *gorm.DB, redis *redis.Client, logger *logrus.Logger) *ScanTaskService {
	return &ScanTaskService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (s *ScanTaskService) GetScanTaskByType(task string) (*dbmodels.ScanTask, error) {
	var scanTask dbmodels.ScanTask
	result := s.db.Where("task_type = ?", task).First(&scanTask)
	if result.Error != nil {
		return nil, result.Error
	}
	return &scanTask, nil
}

func (s *ScanTaskService) UpdateScanTask(task *dbmodels.ScanTask) {
	s.db.Save(task)
}
