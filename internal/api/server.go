package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/xueqianLu/deep-dive-beacon/config"
	"github.com/xueqianLu/deep-dive-beacon/internal/api/handlers"
	"github.com/xueqianLu/deep-dive-beacon/internal/api/middleware"
	"github.com/xueqianLu/deep-dive-beacon/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	config   *config.Config
	db       *gorm.DB
	redis    *redis.Client
	logger   *logrus.Logger
	services *services.Services
	router   *gin.Engine
}

func NewServer(cfg *config.Config, db *gorm.DB, rdb *redis.Client, logger *logrus.Logger) *Server {
	// Initialize services
	svc := services.NewServices(db, rdb, logger, cfg)
	// Initialize handlers
	h := handlers.NewHandlers(svc, logger)

	// Create server instance
	server := &Server{
		config:   cfg,
		db:       db,
		redis:    rdb,
		logger:   logger,
		services: svc,
	}

	// Setup router
	server.setupRouter(h)

	return server
}

func (s *Server) setupRouter(h *handlers.Handlers) {
	// Set Gin mode
	if s.config.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())

	// Health check
	r.GET("/health", h.Health)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", h.Health)
	}

	s.router = r
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	s.logger.Infof("Starting server on %s", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	return server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server...")

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port),
		Handler: s.router,
	}

	return server.Shutdown(ctx)
}
