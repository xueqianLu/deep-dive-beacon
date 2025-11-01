package cmd

import (
	"github.com/xueqianLu/deep-dive-beacon/config"
	"github.com/xueqianLu/deep-dive-beacon/internal/api"
	"github.com/xueqianLu/deep-dive-beacon/internal/database"
	"github.com/xueqianLu/deep-dive-beacon/internal/logger"
	"github.com/xueqianLu/deep-dive-beacon/internal/redis"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the API server",
	Long:  `Start the API server with all endpoints`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg := config.Load()

		// Initialize logger
		log := logger.Init(cfg.Log.Level)

		// Initialize database
		db, err := database.Init(cfg.Database)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}

		// Initialize Redis
		rdb, err := redis.Init(cfg.Redis)
		if err != nil {
			log.Fatalf("Failed to initialize Redis: %v", err)
		}

		// Start API server
		server := api.NewServer(cfg, db, rdb, log)
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations to set up the database schema`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		log := logger.Init(cfg.Log.Level)

		if err := database.Migrate(cfg.Database); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}

		log.Info("Database migrations completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(migrateCmd)

	serverCmd.Flags().StringP("port", "p", "", "Port to run the server on")
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}
