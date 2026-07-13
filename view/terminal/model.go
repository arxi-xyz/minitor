package terminal

import (
	"time"

	"minitor/collector"
	"minitor/config"

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

	wsURL             string
	wsReconnectDelay  time.Duration
	maxProcessLimit   int
	wsCh              chan tea.Msg
}

func InitialModel() Model {
	return NewModel(config.Default())
}

func NewModel(cfg config.Config) Model {
	reconnectDelay, err := cfg.Client.ReconnectDelayDuration()
	if err != nil {
		reconnectDelay = 3 * time.Second
	}

	return Model{
		wsURL:            cfg.Client.WSURL,
		wsReconnectDelay: reconnectDelay,
		maxProcessLimit:  cfg.Socket.MaxProcessLimit,
		wsCh:             make(chan tea.Msg, 64),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		listenWS(m.wsURL, m.wsReconnectDelay, m.maxProcessLimit, m.wsCh),
		waitForWS(m.wsCh),
	)
}
