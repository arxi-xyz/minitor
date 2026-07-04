package collector

import (
	"time"

	"github.com/shirou/gopsutil/mem"
)

type RamMetric struct {
	Total     uint64
	Used      uint64
	Available uint64
	Free      uint64

	Cached  uint64
	Buffers uint64

	Active   uint64
	Inactive uint64

	UsedPercent float64
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

	for {
		select {
		case <-ticker.C:
			ramVal, err := GetRamInfo()
			if err != nil {
				continue
			}
			ch <- ramVal
		}
	}
}
