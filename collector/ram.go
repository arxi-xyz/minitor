package collector

import (
	"time"

	"github.com/shirou/gopsutil/mem"
)

type RamMetric struct {
	Total     uint64 `json:"total"`
	Used      uint64 `json:"used"`
	Available uint64 `json:"available"`
	Free      uint64 `json:"free"`

	Cached  uint64 `json:"cached"`
	Buffers uint64 `json:"buffers"`

	Active   uint64 `json:"active"`
	Inactive uint64 `json:"inactive"`

	UsedPercent float64 `json:"used_percent"`
}

func NewRamMetric() (RamMetric, error) {
	return RamMetric{}, nil
}

func GetRamInfo() (RamMetric, error) {
	ram, err := NewRamMetric()
	if err != nil {
		return RamMetric{}, err
	}

	memoryStat, err := mem.VirtualMemory()
	if err != nil {
		return RamMetric{}, err
	}

	ram.Total = memoryStat.Total
	ram.Used = memoryStat.Used
	ram.Available = memoryStat.Available
	ram.Free = memoryStat.Free

	ram.Cached = memoryStat.Cached
	ram.Buffers = memoryStat.Buffers

	ram.Active = memoryStat.Active
	ram.Inactive = memoryStat.Inactive

	ram.UsedPercent = memoryStat.UsedPercent

	return ram, nil
}

func WorkerRam(ch chan RamMetric) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ramVal, err := GetRamInfo()
		if err != nil {
			continue
		}
		ch <- ramVal
	}
}
