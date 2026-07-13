package terminal

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.processOffset = clampProcessOffset(m.processOffset, m.ProcessMetric, m.processVisibleLines())
		return m, waitForWS(m.wsCh)

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
		return m, waitForWS(m.wsCh)

	case snapshotMsg:
		m.CpuMetric = msg.CPU
		m.RamMetric = msg.RAM
		m.DiskMetric = msg.Disk
		m.NetworkMetric = msg.Network
		return m, waitForWS(m.wsCh)

	case processesMsg:
		m.ProcessMetric = msg.Processes
		m.processOffset = clampProcessOffset(m.processOffset, m.ProcessMetric, m.processVisibleLines())
		return m, waitForWS(m.wsCh)
	}

	return m, waitForWS(m.wsCh)
}
