// collector/cpu.go
package collector

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func GetCpuInfo() (CpuMetric, error) {
	cpu, err := NewCpuMetric()
	if err != nil {
		return CpuMetric{}, err
	}
	return cpu, nil
}

func WorkerCpu(ch chan CpuMetric) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cpuMetric, err := GetCpuInfo()
			if err != nil {
				// Handle error (e.g., log it)
				continue
			}
			ch <- cpuMetric
		}
	}
}

type CpuMetric struct {
	Usage         float64
	CoreUsage     []float64
	LogicalCores  int
	PhysicalCores int
	ModelName     string
}

func NewCpuMetric() (CpuMetric, error) {
	info, err := cpu.Info()
	if err != nil {
		return CpuMetric{}, err
	}

	overallUsage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return CpuMetric{}, err
	}

	coreUsage, err := cpu.Percent(0, true)
	if err != nil {
		return CpuMetric{}, err
	}

	physical, err := cpu.Counts(false)
	if err != nil {
		return CpuMetric{}, err
	}

	logical, err := cpu.Counts(true)
	if err != nil {
		return CpuMetric{}, err
	}

	metric := CpuMetric{
		Usage:         overallUsage[0],
		CoreUsage:     coreUsage,
		PhysicalCores: physical,
		LogicalCores:  logical,
	}

	if len(info) > 0 {
		metric.ModelName = info[0].ModelName
	}

	return metric, nil
}
