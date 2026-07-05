package terminal

import (
	"fmt"
	"minitor/collector"

	"github.com/charmbracelet/lipgloss"
)

const boxWidth = 40
const boxHeight = 22

func (m Model) View() string {
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuBoxView(m.CpuMetric),
		"  ", // spacing
		diskBoxView(m.DiskMetric),
		"  ", // spacing
		ramBoxView(m.RamMetric),
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

func cpuBoxView(cpuMetric collector.CpuMetric) string {
	cores := ""
	for i, v := range cpuMetric.CoreUsage {
		cores += fmt.Sprintf("Core %02d : %5.1f%%\n", i, v)
	}

	return lipgloss.NewStyle().
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
			cpuMetric.ModelName,
			cpuMetric.Usage,
			cpuMetric.PhysicalCores,
			cpuMetric.LogicalCores,
			cores,
		))
}

func diskBoxView(diskMetric collector.DiskMetric) string {
	return lipgloss.NewStyle().
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
			diskMetric.MountPoint,
			diskMetric.FSType,
			diskMetric.UsedPercent,
			diskMetric.Total/1024/1024/1024,
			diskMetric.Free/1024/1024/1024,
			diskMetric.ReadBytes/1024/1024,
			diskMetric.WriteBytes/1024/1024,
		))
}

func ramBoxView(ramMetric collector.RamMetric) string {
	return lipgloss.NewStyle().
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
			ramMetric.UsedPercent,

			ramMetric.Total/1024/1024/1024,
			ramMetric.Used/1024/1024/1024,
			ramMetric.Available/1024/1024/1024,
			ramMetric.Free/1024/1024/1024,

			ramMetric.Cached/1024/1024,
			ramMetric.Buffers/1024/1024,

			ramMetric.Active/1024/1024,
			ramMetric.Inactive/1024/1024,
		))
}
