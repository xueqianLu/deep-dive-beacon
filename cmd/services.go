package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xueqianLu/deep-dive-beacon/config"
	"github.com/xueqianLu/deep-dive-beacon/internal/database"
	"github.com/xueqianLu/deep-dive-beacon/internal/logger"
	"github.com/xueqianLu/deep-dive-beacon/internal/redis"
	"github.com/xueqianLu/deep-dive-beacon/processor/blockscanner"
	"os"
	"os/signal"
	"syscall"
)

var blockScanner = &cobra.Command{
	Use:   "block-scanner",
	Short: "Start the block scanner",
	Long:  `Start the block scanner to sync blockchain data and store to database`,
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

		// Handle graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		scanner := beaconscanner.NewBeaconBlockScanner(cfg, db, rdb, log)

		if err := scanner.Start(); err != nil {
			log.Fatalf("Beacon block scanner failed: %v", err)
		}

		<-sigChan
		log.Info("Received shutdown signal, stopping event processor...")
		scanner.Stop()

		log.Info("Beacon Block scanner stopped")
	},
}

func init() {
	rootCmd.AddCommand(blockScanner)
}
