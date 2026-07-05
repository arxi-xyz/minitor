package terminal

import (
	"minitor/collector"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	CpuMetric  collector.CpuMetric
	RamMetric  collector.RamMetric
	DiskMetric collector.DiskMetric

	cpuChannel  chan collector.CpuMetric
	ramChannel  chan collector.RamMetric
	diskChannel chan collector.DiskMetric
}

func InitialModel() Model {
	cpuChannel := make(chan collector.CpuMetric)
	ramChannel := make(chan collector.RamMetric)
	diskChannel := make(chan collector.DiskMetric)

	return Model{
		cpuChannel:  cpuChannel,
		ramChannel:  ramChannel,
		diskChannel: diskChannel,
	}
}

func (m Model) Init() tea.Cmd {
	go collector.WorkerCpu(m.cpuChannel)
	go collector.WorkerRam(m.ramChannel)
	go collector.WorkerDisk(m.diskChannel)

	return tea.Batch(
		waitForCpuMetric(m.cpuChannel),
		waitForDiskMetric(m.diskChannel),
		waitForRamMetric(m.ramChannel),
	)
}

func waitForRamMetric(ch chan collector.RamMetric) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func waitForDiskMetric(ch chan collector.DiskMetric) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func waitForCpuMetric(ch chan collector.CpuMetric) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
