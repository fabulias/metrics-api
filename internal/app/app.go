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
	// Registrar un nuevo dispositivo
	router.POST("/devices/register", metricsHandler.RegisterDevice)

	// Enviar datos de monitoreo para un dispositivo
	router.POST("/devices/:id/metrics", metricsHandler.SendMetrics)

	// Consultar las últimas métricas de un dispositivo
	router.GET("/devices/:id/metrics", metricsHandler.GetLatestMetrics)

	// Consultar el historial de métricas de un dispositivo
	router.GET("/devices/:id/metrics/history", metricsHandler.GetMetricsHistory)

	return &Application{
		router: router,
		port:   cfg.ServicePort,
	}
}

func (a *Application) Run() {
	a.router.Run(":" + a.port)
}
