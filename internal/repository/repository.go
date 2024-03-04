package repository

import (
	"time"

	metricsapi "github.com/fabulias/metrics-api"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MetricsRepository struct {
	db *gorm.DB
}

func New() *MetricsRepository {
	db, err := gorm.Open(sqlite.Open("db/metrics.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&metricsapi.Device{})
	db.AutoMigrate(&metricsapi.RawMetrics{})
	db.AutoMigrate(&metricsapi.CalculatedMetrics{})
	return &MetricsRepository{db}
}

func (r *MetricsRepository) RegisterDevice(deviceName string) (*metricsapi.Device, error) {
	device := metricsapi.Device{Name: deviceName}
	if result := r.db.Create(&device); result.Error != nil {
		return nil, result.Error
	}
	return &device, nil
}

// func (r *MetricsRepository) SendMetrics(id uint, metrics metricsapi.Metrics) (*metricsapi.Metrics, error) {
// 	metricsResp := metricsapi.Metrics{
// 		DeviceId:    id,
// 		CPUUsage:    metrics.CPUUsage,
// 		MemoryUsage: metrics.MemoryUsage,
// 		DiskUsage:   metrics.DiskUsage,
// 		Timestamp:   metrics.Timestamp,
// 	}

// 	if result := r.db.Create(metricsResp); result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return &metricsResp, nil
// }

func (r *MetricsRepository) GetLatestMetrics(id uint) (*metricsapi.CalculatedMetrics, error) {
	var metric metricsapi.CalculatedMetrics

	// Order by timestamp in descending order
	result := r.db.Where("device_id = ?", id).Order("timestamp desc").First(&metric)

	// Handle errors
	if result.Error != nil {
		return nil, result.Error
	}

	return &metric, nil
}

func (r *MetricsRepository) GetMetricsHistory(id uint) ([]*metricsapi.CalculatedMetrics, error) {
	var metrics []*metricsapi.CalculatedMetrics

	// Order by timestamp in descending order
	result := r.db.Where("device_id =?", id).Order("timestamp desc").Find(&metrics)

	if result.Error != nil {
		return nil, result.Error
	}

	return metrics, nil
}

func (r *MetricsRepository) DeviceExists(id uint) bool {
	if result := r.db.Model(&metricsapi.Device{}).Where("id =?", id).First(&metricsapi.Device{}); result.Error != nil {
		return false
	}
	return true
}

func (r *MetricsRepository) GetCalculatedMetrics(deviceID uint, timestamp time.Time) (*metricsapi.CalculatedMetrics, error) {
	var cm metricsapi.CalculatedMetrics
	result := r.db.Where("device_id = ?", deviceID).Order("timestamp desc").First(&cm)
	if result.Error != nil {
		return nil, result.Error
	}

	return &cm, nil
}

func (r *MetricsRepository) CreateRawMetrics(tx *gorm.DB, rm *metricsapi.RawMetrics) error {
	return tx.Create(rm).Error
}

func (r *MetricsRepository) CreateCalculatedMetrics(tx *gorm.DB, cm *metricsapi.CalculatedMetrics) error {
	return tx.Save(cm).Error
}

func (r *MetricsRepository) Begin() *gorm.DB {
	return r.db.Begin()
}
