package service

import (
	"errors"
	"math"

	metricsapi "github.com/fabulias/metrics-api"
	"github.com/fabulias/metrics-api/internal/repository"
	"gorm.io/gorm"
)

type MetricsService struct {
	repository *repository.MetricsRepository
}

func New(repository *repository.MetricsRepository) *MetricsService {
	return &MetricsService{
		repository: repository,
	}
}

func (s *MetricsService) RegisterDevice(deviceName string) (*metricsapi.Device, error) {
	return s.repository.RegisterDevice(deviceName)
}
func (s *MetricsService) SendMetrics(deviceID uint, rm *metricsapi.RawMetrics) error {
	// Check if device exists
	if !s.repository.DeviceExists(deviceID) {
		return errors.New("device not found")
	}

	// Obtain previous metrics calculated for this device
	cm, err := s.repository.GetCalculatedMetrics(deviceID, rm.Timestamp)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	rm.DeviceID = deviceID
	// Calculate new metrics
	newAvg, newMax, newMin := calculateMetrics(rm, cm)

	// If there aren't preexistent metrics, create new ones.
	// if cm == nil {
	cm = &metricsapi.CalculatedMetrics{
		DeviceID:  deviceID,
		Timestamp: rm.Timestamp,
		CPUAvg:    newAvg.CPU,
		CPUMax:    newMax.CPU,
		CPUMin:    newMin.CPU,
		MemoryAvg: newAvg.Memory,
		MemoryMax: newMax.Memory,
		MemoryMin: newMin.Memory,
		DiskAvg:   newAvg.Disk,
		DiskMax:   newMax.Disk,
		DiskMin:   newMin.Disk,
	}
	// } else {
	// 	// Update current metrics
	// 	cm.CPUAvg = newAvg.CPU
	// 	cm.CPUMax = newMax.CPU
	// 	cm.CPUMin = newMin.CPU
	// 	cm.MemoryAvg = newAvg.Memory
	// 	cm.MemoryMax = newMax.Memory
	// 	cm.MemoryMin = newMin.Memory
	// 	cm.DiskAvg = newAvg.Disk
	// 	cm.DiskMax = newMax.Disk
	// 	cm.DiskMin = newMin.Disk
	// }

	// Init transaction
	tx := s.repository.Begin()

	// Store raw metrics
	err = s.repository.CreateRawMetrics(tx, rm)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Save new calculated metrics
	err = s.repository.CreateCalculatedMetrics(tx, cm)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Confirm transaction
	tx.Commit()

	return nil
}

func (s *MetricsService) GetLatestMetrics(id uint) (*metricsapi.CalculatedMetrics, error) {
	return s.repository.GetLatestMetrics(id)
}

func (s *MetricsService) GetMetricsHistory(id uint) ([]*metricsapi.CalculatedMetrics, error) {
	return s.repository.GetMetricsHistory(id)
}

type Metrics struct {
	CPU    float64
	Memory float64
	Disk   float64
}

func calculateMetrics(rm *metricsapi.RawMetrics, cm *metricsapi.CalculatedMetrics) (Metrics, Metrics, Metrics) {
	var newAvg, newMax, newMin Metrics

	// If existing calculated metrics are not available, use the current raw value as initial values
	if cm == nil {
		newAvg = Metrics{CPU: rm.CPUUsage, Memory: rm.MemoryUsage, Disk: rm.DiskUsage}
		newMax = newAvg
		newMin = newAvg
	} else {
		// Calculate average
		newAvg.CPU = (cm.CPUAvg + rm.CPUUsage) / 2
		newAvg.Memory = (cm.MemoryAvg + rm.MemoryUsage) / 2
		newAvg.Disk = (cm.DiskAvg + rm.DiskUsage) / 2

		// Update existing max and min values
		newMax.CPU = math.Max(cm.CPUMax, rm.CPUUsage)
		newMax.Memory = math.Max(cm.MemoryMax, rm.MemoryUsage)
		newMax.Disk = math.Max(cm.DiskMax, rm.DiskUsage)

		newMin.CPU = math.Min(cm.CPUMin, rm.CPUUsage)
		newMin.Memory = math.Min(cm.MemoryMin, rm.MemoryUsage)
		newMin.Disk = math.Min(cm.DiskMin, rm.DiskUsage)
	}

	return newAvg, newMax, newMin
}
