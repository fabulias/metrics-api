package app

import (
	"github.com/fabulias/metrics-api/internal/config"
	"github.com/fabulias/metrics-api/internal/handlers"
	"github.com/fabulias/metrics-api/internal/repository"
	"github.com/fabulias/metrics-api/internal/service"

	"github.com/gin-gonic/gin"
)

type Application struct {
	router *gin.Engine
	port   string
}

func NewApplication(cfg config.Config) *Application {
	router := gin.Default()

	repository := repository.New()
	service := service.New(repository)
	metricsHandler := handlers.New(service)
	// Register a new device
	router.POST("/devices/register", metricsHandler.RegisterDevice)

	// Send data of a device related to metrics
	router.POST("/devices/:id/metrics", metricsHandler.SendMetrics)

	// Request latest metrics for a device
	router.GET("/devices/:id/metrics", metricsHandler.GetLatestMetrics)

	// Request history of metrics for a device
	router.GET("/devices/:id/metrics/history", metricsHandler.GetMetricsHistory)

	return &Application{
		router: router,
		port:   cfg.ServicePort,
	}
}

func (a *Application) Run() {
	a.router.Run(":" + a.port)
}
