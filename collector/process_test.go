package collector

import (
	"testing"
	"time"

	pr "github.com/shirou/gopsutil/v4/process"
)

func TestGetProccess(t *testing.T) {
	rawProcesses, _ := pr.Processes()

	for _, p := range rawProcesses {
		_, _ = p.CPUPercent()
	}
	time.Sleep(100 * time.Millisecond)

	procMap := make(map[int]*ProcessMetric)

	for _, p := range rawProcesses {
		name, err := p.Name()

		if err != nil {
			t.Fatal(err)
		}

		mem, err := p.MemoryPercent()
		if err != nil {
			t.Fatal(err)
		}

		cpu, err := p.CPUPercent()
		if err != nil {
			t.Fatal(err)
		}

		status, err := p.Status()
		if err != nil {
			t.Fatal(err)
		}

		user, err := p.Username()

		if err != nil {
			continue
		}

		ppid, err := p.Ppid()

		if err != nil {
			t.Fatal(err)
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
}
