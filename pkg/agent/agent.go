package agent

import (
	"fmt"
	"time"

	metricsapi "github.com/fabulias/metrics-api"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// Agent struct for encapsulating metrics-related methods
type Agent struct{}

// NewAgent creates a new Agent instance with a specified interval
func NewAgent() *Agent {
	return &Agent{}
}

// CollectMetrics concurrently collects system metrics and aggregates results
func (a *Agent) CollectMetrics() (metricsapi.RawMetrics, error) {
	var metrics metricsapi.RawMetrics

	// Create channels for each metric type
	cpuCh := make(chan float64, 1)
	memCh := make(chan float64, 1)
	diskCh := make(chan float64, 1)

	// Launch concurrent goroutines for each metric retrieval
	go func() {
		cpuUsage, err := cpu.Percent(time.Second, false)
		if err != nil {
			fmt.Println("Error obtaining CPU usage:", err)
			cpuCh <- -1 // Signal error
			return
		}
		cpuCh <- cpuUsage[0]
	}()

	go func() {
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Error obtaining memory usage:", err)
			memCh <- -1 // Signal error
			return
		}
		memCh <- memInfo.UsedPercent
	}()

	go func() {
		diskInfo, err := disk.Usage("/")
		if err != nil {
			fmt.Println("Error obtaining disk usage:", err)
			diskCh <- -1 // Signal error
			return
		}
		diskCh <- diskInfo.UsedPercent
	}()

	// Wait for all goroutines to finish
	cpuUsage := <-cpuCh
	if cpuUsage == -1 {
		return metrics, fmt.Errorf("failed to obtain CPU usage")
	}

	memoryUsage := <-memCh
	if memoryUsage == -1 {
		return metrics, fmt.Errorf("failed to obtain memory usage")
	}

	diskUsage := <-diskCh

	// Aggregate and set metrics
	metrics.CPUUsage = cpuUsage
	metrics.MemoryUsage = memoryUsage
	metrics.DiskUsage = diskUsage
	metrics.Timestamp = time.Now()

	return metrics, nil
}
