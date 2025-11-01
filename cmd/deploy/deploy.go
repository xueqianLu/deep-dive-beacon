package deploy

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/xueqianLu/deep-dive-beacon/cmd/deploy/production"
	"github.com/xueqianLu/deep-dive-beacon/config"
	"github.com/xueqianLu/deep-dive-beacon/types"
	"gorm.io/gorm"
	"os"
)

type DeployInstance interface {
	Execute() error
}

func GetDeployInstance(db *gorm.DB, redis *redis.Client, logger *logrus.Logger, cfg *config.Config, depcfg types.DeployConfig) DeployInstance {
	return production.NewProdDeploy(db, redis, logger, cfg, depcfg)
}

func ParseDeployConfig(depConfig string) (types.DeployConfig, error) {
	var depcfg types.DeployConfig
	if data, err := os.ReadFile(depConfig); err != nil {
		return depcfg, err
	} else {
		if err = json.Unmarshal(data, &depcfg); err != nil {
			return depcfg, err
		}
		return depcfg, nil
	}
}
