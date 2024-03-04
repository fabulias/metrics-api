package handlers

import (
	"net/http"
	"strconv"

	metricsapi "github.com/fabulias/metrics-api"
	"github.com/fabulias/metrics-api/internal/service"
	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	service *service.MetricsService
}

func New(svc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		service: svc,
	}
}

func (h *MetricsHandler) RegisterDevice(c *gin.Context) {

	type Device struct {
		Name string `json:"name"`
	}

	var device Device
	if err := c.BindJSON(&device); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	deviceResp, err := h.service.RegisterDevice(device.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, &deviceResp)
}

func (h *MetricsHandler) SendMetrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var metric metricsapi.RawMetrics
	if err := c.BindJSON(&metric); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SendMetrics(uint(id), &metric)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Metrics created successfully"})
}

func (h *MetricsHandler) GetLatestMetrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	latestMetrics, err := h.service.GetLatestMetrics(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": latestMetrics})

}

func (h *MetricsHandler) GetMetricsHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metricsHistory, err := h.service.GetMetricsHistory(uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": metricsHistory})
}
