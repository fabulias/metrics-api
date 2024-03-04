package metricsapi

import (
	"time"

	"gorm.io/gorm"
)

type RawMetrics struct {
	gorm.Model
	DeviceID    uint      `json:"deviceId" gorm:"foreignKey:ID;references:devices"`
	CPUUsage    float64   `json:"cpuUsage"`
	MemoryUsage float64   `json:"memoryUsage"`
	DiskUsage   float64   `json:"diskUsage"`
	Timestamp   time.Time `json:"timestamp" gorm:"index"`
}

type CalculatedMetrics struct {
	gorm.Model
	DeviceID  uint      `json:"deviceId" gorm:"foreignKey:ID;references:devices"`
	Timestamp time.Time `json:"timestamp" gorm:"index"`
	CPUAvg    float64   `json:"cpuAvg"`
	CPUMax    float64   `json:"cpuMax"`
	CPUMin    float64   `json:"cpuMin"`
	MemoryAvg float64   `json:"memoryAvg"`
	MemoryMax float64   `json:"memoryMax"`
	MemoryMin float64   `json:"memoryMin"`
	DiskAvg   float64   `json:"diskAvg"`
	DiskMax   float64   `json:"diskMax"`
	DiskMin   float64   `json:"diskMin"`
}

type Device struct {
	gorm.Model
	Name string `json:"name"`
}
