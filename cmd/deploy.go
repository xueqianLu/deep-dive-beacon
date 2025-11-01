package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xueqianLu/deep-dive-beacon/cmd/deploy"
	"github.com/xueqianLu/deep-dive-beacon/config"
	"github.com/xueqianLu/deep-dive-beacon/internal/database"
	"github.com/xueqianLu/deep-dive-beacon/internal/logger"
	"github.com/xueqianLu/deep-dive-beacon/internal/redis"
	"strings"
)

var (
	depConfig string
	migratedb bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Initialize database and add initial data",
	Long:  `Run database migrations and insert initial data (default chain, token, admin user, etc.)`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		log := logger.Init(cfg.Log.Level)

		log.Info("Initializing database...")
		if migratedb {
			if err := database.Migrate(cfg.Database); err != nil {
				log.Fatalf("Database migration failed: %v", err)
			}
		}

		// check database already initialized
		db, err := database.Init(cfg.Database)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		// Parse deployment config
		deployCfg, err := deploy.ParseDeployConfig(depConfig)
		if err != nil {
			log.Fatalf("Failed to parse deployment config: %v", err)
		}

		// Initialize Redis
		rdb, err := redis.Init(cfg.Redis)
		if err != nil {
			log.Fatalf("Failed to initialize Redis: %v", err)
		}

		log.WithField("config", deployCfg).Info("Inserting initial data...")

		// Begin transaction
		tx := db.Begin()
		if tx.Error != nil {
			log.Fatalf("Failed to begin transaction: %v", tx.Error)
		}

		instance := deploy.GetDeployInstance(tx, rdb, log, cfg, deployCfg)
		if err := instance.Execute(); err != nil {
			tx.Rollback()
			log.WithError(err).Fatal("Deployment failed, please check error and deploy again.")
		} else {
			if err := tx.Commit().Error; err != nil {
				log.WithError(err).Fatalf("Deployment failed, commit transaction to db failed")
			}
			log.Info("Deployment succeeded.")
		}

	},
}

func splitAndTrim(s string) []string {
	items := []string{}
	for _, item := range strings.Split(s, ",") {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func init() {
	deployCmd.Flags().StringVar(&depConfig, "deploy-config", "deploy.json", "Deployment configuration file path")
	deployCmd.Flags().BoolVar(&migratedb, "migrate-db", true, "Run database migrations before deployment")
	rootCmd.AddCommand(deployCmd)
}
