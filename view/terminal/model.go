package terminal

import (
	"minitor/collector"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	CpuMetric     collector.CpuMetric
	RamMetric     collector.RamMetric
	DiskMetric    collector.DiskMetric
	NetworkMetric collector.NetworkMetric
	ProcessMetric []*collector.ProcessMetric

	width         int
	height        int
	processOffset int

	cpuChannel     chan collector.CpuMetric
	ramChannel     chan collector.RamMetric
	diskChannel    chan collector.DiskMetric
	networkChannel chan collector.NetworkMetric
	processChannel chan []*collector.ProcessMetric
}

func InitialModel() Model {
	cpuChannel := make(chan collector.CpuMetric)
	ramChannel := make(chan collector.RamMetric)
	diskChannel := make(chan collector.DiskMetric)
	networkChannel := make(chan collector.NetworkMetric)
	processChannel := make(chan []*collector.ProcessMetric)

	return Model{
		cpuChannel:     cpuChannel,
		ramChannel:     ramChannel,
		diskChannel:    diskChannel,
		networkChannel: networkChannel,
		processChannel: processChannel,
	}
}

func (m Model) Init() tea.Cmd {
	go collector.WorkerCpu(m.cpuChannel)
	go collector.WorkerRam(m.ramChannel)
	go collector.WorkerDisk(m.diskChannel)
	go collector.WorkerNetwork(m.networkChannel)
	go collector.ProcessWorker(m.processChannel)

	return tea.Batch(
		waitForCpuMetric(m.cpuChannel),
		waitForDiskMetric(m.diskChannel),
		waitForRamMetric(m.ramChannel),
		waitForNetworkMetric(m.networkChannel),
		waitForProcessMetric(m.processChannel),
	)
}

func waitForRamMetric(ch chan collector.RamMetric) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}

func waitForNetworkMetric(ch chan collector.NetworkMetric) tea.Cmd {
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

func waitForProcessMetric(ch chan []*collector.ProcessMetric) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
