package terminal

import (
	"minitor/collector"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case collector.RamMetric:
		m.RamMetric = msg
		return m, waitForRamMetric(m.ramChannel)

	case collector.CpuMetric:
		m.CpuMetric = msg
		return m, waitForCpuMetric(m.cpuChannel)
	case collector.DiskMetric:
		m.DiskMetric = msg
		return m, waitForDiskMetric(m.diskChannel)

	case collector.NetworkMetric:
		m.NetworkMetric = msg
		return m, waitForNetworkMetric(m.networkChannel)
	case []*collector.ProcessMetric:
		m.ProcessMetric = msg
		return m, waitForNetworkMetric(m.networkChannel)
	}
	return m, nil
}
