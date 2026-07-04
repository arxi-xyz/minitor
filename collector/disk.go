package collector

import (
	"time"

	"github.com/shirou/gopsutil/disk"
)

type DiskMetric struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64

	ReadBytes  uint64
	WriteBytes uint64

	MountPoint string
	FSType     string
}

func NewDiskMetric() DiskMetric {
	return DiskMetric{}
}

func GetDiskInfo() (DiskMetric, error) {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return DiskMetric{}, err
	}

	io, _ := disk.IOCounters()

	var read, write uint64
	for _, v := range io {
		read += v.ReadBytes
		write += v.WriteBytes
	}

	return DiskMetric{
		Total:       diskUsage.Total,
		Used:        diskUsage.Used,
		Free:        diskUsage.Free,
		UsedPercent: diskUsage.UsedPercent,
		MountPoint:  "/",
		FSType:      diskUsage.Fstype,
		ReadBytes:   read,
		WriteBytes:  write,
	}, nil
}

func WorkerDisk(ch chan<- DiskMetric) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		diskMetric, err := GetDiskInfo()
		if err != nil {
			// Handle error (e.g., log it)
			return
		}
		ch <- diskMetric
	}
}
