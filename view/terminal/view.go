package terminal

import (
	"fmt"
	"minitor/collector"
	"minitor/helper"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const boxWidth = 30
const boxHeight = 22

func (m Model) View() string {
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpuBoxView(m.CpuMetric),
		"  ",
		diskBoxView(m.DiskMetric),
		"  ",
		ramBoxView(m.RamMetric),
		"  ",
		networkBoxView(m.NetworkMetric),
	)

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Render("MINITOR - SYSTEM DASHBOARD")

	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		Render(strings.Repeat("─", 132))

	processSection := processListView(m.ProcessMetric, m.processOffset, m.processVisibleLines())

	footer := lipgloss.NewStyle().
		Faint(true).
		Render("j/k scroll processes · q quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		row,
		"",
		separator,
		"",
		processSection,
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

const processSectionOverhead = 34

func (m Model) processVisibleLines() int {
	if m.height <= processSectionOverhead {
		return 10
	}
	return m.height - processSectionOverhead
}

func processListView(roots []*collector.ProcessMetric, offset, maxLines int) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Render("PROCESSES")

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Render(fmt.Sprintf("%-8s %6s %6s %-12s %s", "PID", "CPU%", "MEM%", "USER", "NAME"))

	if len(roots) == 0 {
		return title + "\n\n" + header + "\n\n" + lipgloss.NewStyle().Faint(true).Render("loading processes...")
	}

	allLines := collectProcessLines(roots)
	offset = clampProcessOffset(offset, roots, maxLines)

	end := offset + maxLines
	if end > len(allLines) {
		end = len(allLines)
	}

	visible := allLines[offset:end]

	scrollInfo := ""
	if len(allLines) > maxLines {
		scrollInfo = lipgloss.NewStyle().Faint(true).Render(
			fmt.Sprintf("  %d–%d of %d", offset+1, end, len(allLines)),
		)
	}

	return title + scrollInfo + "\n\n" + header + "\n" + strings.Join(visible, "\n")
}

func collectProcessLines(roots []*collector.ProcessMetric) []string {
	var lines []string
	for _, proc := range sortedProcesses(roots) {
		lines = append(lines, renderProcessTree(proc, 0)...)
	}
	return lines
}

func maxProcessOffset(roots []*collector.ProcessMetric, maxLines int) int {
	total := len(collectProcessLines(roots))
	if total <= maxLines {
		return 0
	}
	return total - maxLines
}

func clampProcessOffset(offset int, roots []*collector.ProcessMetric, maxLines int) int {
	maxOffset := maxProcessOffset(roots, maxLines)
	if offset < 0 {
		return 0
	}
	if offset > maxOffset {
		return maxOffset
	}
	return offset
}

func renderProcessTree(proc *collector.ProcessMetric, depth int) []string {
	prefix := treePrefix(depth)

	line := fmt.Sprintf(
		"%s%-8d %5.1f%% %5.1f%% %-12s %s",
		prefix,
		proc.PID,
		proc.Cpu,
		proc.Mem,
		truncate(proc.User, 12),
		proc.Name,
	)

	lines := []string{line}

	for _, child := range sortedChildren(proc) {
		lines = append(lines, renderProcessTree(child, depth+1)...)
	}

	return lines
}

func treePrefix(depth int) string {
	if depth == 0 {
		return ""
	}
	return strings.Repeat("  ", depth-1) + "└─ "
}

func sortedProcesses(procs []*collector.ProcessMetric) []*collector.ProcessMetric {
	sorted := make([]*collector.ProcessMetric, len(procs))
	copy(sorted, procs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].PID < sorted[j].PID
	})
	return sorted
}

func sortedChildren(proc *collector.ProcessMetric) []*collector.ProcessMetric {
	children := make([]*collector.ProcessMetric, 0, len(proc.Childs))
	for _, c := range proc.Childs {
		children = append(children, c)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].PID < children[j].PID
	})
	return children
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
