package terminal

import (
	"minitor/collector"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.processOffset = clampProcessOffset(m.processOffset, m.ProcessMetric, m.processVisibleLines())
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.processOffset > 0 {
				m.processOffset--
			}
		case "down", "j":
			maxOffset := maxProcessOffset(m.ProcessMetric, m.processVisibleLines())
			if m.processOffset < maxOffset {
				m.processOffset++
			}
		case "pgup":
			m.processOffset -= 10
			if m.processOffset < 0 {
				m.processOffset = 0
			}
		case "pgdown":
			m.processOffset += 10
			maxOffset := maxProcessOffset(m.ProcessMetric, m.processVisibleLines())
			if m.processOffset > maxOffset {
				m.processOffset = maxOffset
			}
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
		m.processOffset = clampProcessOffset(m.processOffset, m.ProcessMetric, m.processVisibleLines())
		return m, waitForProcessMetric(m.processChannel)
	}
	return m, nil
}
