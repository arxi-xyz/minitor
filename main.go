package main

import (
	"fmt"
	"os"

	"minitor/collector"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const boxWidth = 40
const boxHeight = 22

type Model struct {
	CpuMetric  collector.CpuMetric
	RamMetric  collector.RamMetric
	DiskMetric collector.DiskMetric

	cpuChannel  chan collector.CpuMetric
	ramChannel  chan collector.RamMetric
	diskChannel chan collector.DiskMetric
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
func initialModel() Model {
	cpuChannel := make(chan collector.CpuMetric)
	ramChannel := make(chan collector.RamMetric)
	diskChannel := make(chan collector.DiskMetric)

	return Model{
		cpuChannel:  cpuChannel,
		ramChannel:  ramChannel,
		diskChannel: diskChannel,
	}
}

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
	}
	return m, nil
}

func (m Model) View() string {
	// ---------- CPU cores ----------
	cores := ""
	for i, v := range m.CpuMetric.CoreUsage {
		cores += fmt.Sprintf("Core %02d : %5.1f%%\n", i, v)
	}

	cpuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1).
		Width(boxWidth).
		Height(boxHeight).
		Render(fmt.Sprintf(
			`CPU
-------------------
Model  : %s
Usage  : %.1f%%
Phys   : %d
Logic  : %d

%s`,
			m.CpuMetric.ModelName,
			m.CpuMetric.Usage,
			m.CpuMetric.PhysicalCores,
			m.CpuMetric.LogicalCores,
			cores,
		))

	// ---------- DISK ----------
	diskBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#43BF6D")).
		Padding(1).
		Width(boxWidth).
		Height(boxHeight).
		Render(fmt.Sprintf(
			`DISK
-------------------
Mount : %s
Type  : %s

Used  : %.1f%%
Total : %d GB
Free  : %d GB

Read  : %d MB
Write : %d MB`,
			m.DiskMetric.MountPoint,
			m.DiskMetric.FSType,
			m.DiskMetric.UsedPercent,
			m.DiskMetric.Total/1024/1024/1024,
			m.DiskMetric.Free/1024/1024/1024,
			m.DiskMetric.ReadBytes/1024/1024,
			m.DiskMetric.WriteBytes/1024/1024,
		))

	// ---------- RAM ----------
	ramBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E2B714")).
		Padding(1).
		Width(boxWidth).
		Height(boxHeight).
		Render(fmt.Sprintf(
			`RAM
-------------------
Usage      : %.1f%%

Total      : %d GB
Used       : %d GB
Available  : %d GB
Free       : %d GB

Cached     : %d MB
Buffers    : %d MB

Active     : %d MB
Inactive   : %d MB`,
			m.RamMetric.UsedPercent,

			m.RamMetric.Total/1024/1024/1024,
			m.RamMetric.Used/1024/1024/1024,
			m.RamMetric.Available/1024/1024/1024,
			m.RamMetric.Free/1024/1024/1024,

			m.RamMetric.Cached/1024/1024,
			m.RamMetric.Buffers/1024/1024,

			m.RamMetric.Active/1024/1024,
			m.RamMetric.Inactive/1024/1024,
		))

	// ---------- LAYOUT ----------
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuBox,
		"  ", // spacing
		diskBox,
		"  ", // spacing
		ramBox,
	)

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Render("MINITOR - SYSTEM DASHBOARD")

	footer := lipgloss.NewStyle().
		Faint(true).
		Render("press q to quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		row,
		"",
		footer,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}
