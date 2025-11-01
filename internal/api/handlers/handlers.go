package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/xueqianLu/deep-dive-beacon/internal/services"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Handlers struct {
	services *services.Services
	logger   *logrus.Logger
}

func NewHandlers(services *services.Services, logger *logrus.Logger) *Handlers {
	return &Handlers{
		services: services,
		logger:   logger,
	}
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "wallet-service",
	})
}
