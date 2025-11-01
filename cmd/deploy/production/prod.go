package production

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/deep-dive-beacon/config"
	"github.com/xueqianLu/deep-dive-beacon/types"
	"gorm.io/gorm"
)

type ProdDeploy struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *logrus.Logger
	cfg    *config.Config
	depcfg types.DeployConfig
}

func NewProdDeploy(db *gorm.DB, redis *redis.Client, logger *logrus.Logger, cfg *config.Config, depcfg types.DeployConfig) ProdDeploy {
	return ProdDeploy{
		db:     db,
		redis:  redis,
		logger: logger,
		cfg:    cfg,
		depcfg: depcfg,
	}
}

func (d ProdDeploy) Execute() error {
	return nil
}
