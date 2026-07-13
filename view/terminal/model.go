package terminal

import (
	"minitor/collector"

	tea "github.com/charmbracelet/bubbletea"
)

const defaultWSURL = "ws://127.0.0.1:8080/ws"

type Model struct {
	CpuMetric     collector.CpuMetric
	RamMetric     collector.RamMetric
	DiskMetric    collector.DiskMetric
	NetworkMetric collector.NetworkMetric
	ProcessMetric []*collector.ProcessMetric

	width         int
	height        int
	processOffset int

	wsURL string
	wsCh  chan tea.Msg
}

func InitialModel() Model {
	return InitialModelWithURL(defaultWSURL)
}

func InitialModelWithURL(url string) Model {
	return Model{
		wsURL: url,
		wsCh:  make(chan tea.Msg, 64),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		listenWS(m.wsURL, m.wsCh),
		waitForWS(m.wsCh),
	)
}
