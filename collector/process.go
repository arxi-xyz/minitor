package collector

import (
	"time"

	pr "github.com/shirou/gopsutil/v4/process"
)

type ProcessMetric struct {
	PID    int
	PPID   int
	Name   string
	Mem    float64
	Cpu    float64
	Status string
	User   string
	Childs map[int]*ProcessMetric
}

func GetProccessInfo() ([]*ProcessMetric, error) {
	rawProcesses, _ := pr.Processes()

	for _, p := range rawProcesses {
		_, _ = p.CPUPercent()
	}
	time.Sleep(100 * time.Millisecond)

	procMap := make(map[int]*ProcessMetric)

	var roots []*ProcessMetric

	for _, p := range rawProcesses {
		name, err := p.Name()

		if err != nil {
			return roots, err
		}

		mem, err := p.MemoryPercent()
		if err != nil {
			return roots, err
		}

		cpu, err := p.CPUPercent()
		if err != nil {
			return roots, err
		}

		status, err := p.Status()
		if err != nil {
			return roots, err
		}

		user, err := p.Username()

		if err != nil {
			continue
		}

		ppid, err := p.Ppid()

		if err != nil {
			return roots, err
		}
		Process := ProcessMetric{
			Name:   name,
			PPID:   int(ppid),
			Mem:    float64(mem),
			Cpu:    float64(cpu),
			Status: status[0],
			User:   user,
			PID:    int(p.Pid),
			Childs: make(map[int]*ProcessMetric),
		}

		procMap[int(p.Pid)] = &Process
	}

	for _, proc := range procMap {
		parent, ok := procMap[int(proc.PPID)]

		if ok {
			parent.Childs[int(proc.PID)] = proc
		} else {
			roots = append(roots, proc)
		}
	}
	return roots, nil
}

func ProcessWorker(ch chan<- []*ProcessMetric) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metric, err := GetProccessInfo()
		if err != nil {
			continue
		}

		ch <- metric
	}
}
