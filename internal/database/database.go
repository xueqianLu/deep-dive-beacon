package database

import (
	"fmt"
	"github.com/xueqianLu/deep-dive-beacon/models/dbmodels"
	"log"
	"os"
	"time"

	"github.com/xueqianLu/deep-dive-beacon/config"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			IgnoreRecordNotFoundError: true,
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info, // Silent, Error, Warn, Info
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(cfg config.DatabaseConfig) error {
	db, err := Init(cfg)
	if err != nil {
		return err
	}

	// Auto migrate all models
	return db.AutoMigrate(
		&dbmodels.BeaconBlock{},
		&dbmodels.BeaconBlock{},
		&dbmodels.ScanTask{},
		&dbmodels.DirectlyScanTask{},
		&dbmodels.BeaconAttestation{},
	)
}
