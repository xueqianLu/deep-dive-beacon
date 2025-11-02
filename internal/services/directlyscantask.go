package services

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/deep-dive-beacon/models/dbmodels"
	"gorm.io/gorm"
)

type DirectlyScanTaskService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
}

func NewDirectlyScanTaskService(db *gorm.DB, redis *redis.Client, logger *logrus.Logger) *DirectlyScanTaskService {
	return &DirectlyScanTaskService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (s *DirectlyScanTaskService) GetScanTaskByType(task string) ([]*dbmodels.DirectlyScanTask, error) {
	var scanTasks []*dbmodels.DirectlyScanTask
	result := s.db.Where("task_type = ?", task).Find(&scanTasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return scanTasks, nil
}

func (s *DirectlyScanTaskService) UpdateScanTask(task *dbmodels.DirectlyScanTask) {
	s.db.Save(task)
}
