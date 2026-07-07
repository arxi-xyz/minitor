package terminal

import (
	"fmt"
	"minitor/collector"
	"minitor/helper"

	"github.com/charmbracelet/lipgloss"
)

const boxWidth = 30
const boxHeight = 22

func (m Model) View() string {
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuBoxView(m.CpuMetric),
		"  ", // spacing
		diskBoxView(m.DiskMetric),
		"  ", // spacing
		ramBoxView(m.RamMetric),
		"  ",
		networkBoxView(m.NetworkMetric),
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
Total : %.2f GB
Free  : %.2f GB

Read  : %.2f MB
Write : %.2f MB`,
			diskMetric.MountPoint,
			diskMetric.FSType,
			diskMetric.UsedPercent,
			helper.Convert(helper.Byte, helper.GigaByte, float64(diskMetric.Total)),
			helper.Convert(helper.Byte, helper.GigaByte, float64(diskMetric.Free)),
			helper.Convert(helper.Byte, helper.MegaByte, float64(diskMetric.ReadBytes)),
			helper.Convert(helper.Byte, helper.MegaByte, float64(diskMetric.WriteBytes)),
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

Total      : %.2f GB
Used       : %.2f GB
Available  : %.2f GB
Free       : %.2f GB

Cached     : %.2f MB
Buffers    : %.2f MB

Active     : %.2f MB
Inactive   : %.2f MB`,
			ramMetric.UsedPercent,

			helper.Convert(helper.Byte, helper.GigaByte, float64(ramMetric.Total)),
			helper.Convert(helper.Byte, helper.GigaByte, float64(ramMetric.Used)),
			helper.Convert(helper.Byte, helper.GigaByte, float64(ramMetric.Available)),
			helper.Convert(helper.Byte, helper.GigaByte, float64(ramMetric.Free)),

			helper.Convert(helper.Byte, helper.MegaByte, float64(ramMetric.Cached)),
			helper.Convert(helper.Byte, helper.MegaByte, float64(ramMetric.Buffers)),

			helper.Convert(helper.Byte, helper.MegaByte, float64(ramMetric.Active)),
			helper.Convert(helper.Byte, helper.MegaByte, float64(ramMetric.Inactive)),
		))
}

func networkBoxView(networkMetric collector.NetworkMetric) string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#E2B714")).
		Padding(1).
		Width(boxWidth).
		Height(boxHeight).
		Render(fmt.Sprintf(
			`NETWORK
-------------------
Interface : %s
IP        : %s

RX Total : %.2f GB
TX Total : %.2f GB

Packets RX : %d
Packets TX : %d

Errors : %d
Drops  : %d`,
			networkMetric.InterfaceName,
			networkMetric.IP,

			helper.Convert(
				helper.Byte,
				helper.GigaByte,
				float64(networkMetric.RxBytes),
			),

			helper.Convert(
				helper.Byte,
				helper.GigaByte,
				float64(networkMetric.TxBytes),
			),

			networkMetric.PacketsRx,
			networkMetric.PacketsTx,

			networkMetric.Errors,
			networkMetric.Drops,
		))
}
