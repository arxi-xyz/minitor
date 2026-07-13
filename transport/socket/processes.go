package socket

import (
	"encoding/json"
	"sort"

	"minitor/collector"
)

const (
	defaultProcessLimit = 50
	maxProcessLimit     = 200
)

type ProcessRow struct {
	PID    int     `json:"pid"`
	PPID   int     `json:"ppid"`
	Name   string  `json:"name"`
	Mem    float64 `json:"mem"`
	Cpu    float64 `json:"cpu"`
	Status string  `json:"status"`
	User   string  `json:"user"`
	Depth  int     `json:"depth"`
}

type ProcessesRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ProcessesPage struct {
	Offset int          `json:"offset"`
	Limit  int          `json:"limit"`
	Total  int          `json:"total"`
	Items  []ProcessRow `json:"items"`
}

func (p ProcessesPage) Message() (Message, error) {
	data, err := json.Marshal(p)
	if err != nil {
		return Message{}, err
	}

	return Message{Type: "processes", Data: data}, nil
}

func flattenProcesses(roots []*collector.ProcessMetric) []ProcessRow {
	rows := make([]ProcessRow, 0, len(roots))

	var walk func(proc *collector.ProcessMetric, depth int)
	walk = func(proc *collector.ProcessMetric, depth int) {
		rows = append(rows, ProcessRow{
			PID:    proc.PID,
			PPID:   proc.PPID,
			Name:   proc.Name,
			Mem:    proc.Mem,
			Cpu:    proc.Cpu,
			Status: proc.Status,
			User:   proc.User,
			Depth:  depth,
		})

		for _, child := range sortedProcessChildren(proc) {
			walk(child, depth+1)
		}
	}

	for _, root := range sortedProcessRoots(roots) {
		walk(root, 0)
	}

	return rows
}

func sortedProcessRoots(procs []*collector.ProcessMetric) []*collector.ProcessMetric {
	sorted := make([]*collector.ProcessMetric, len(procs))
	copy(sorted, procs)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].PID < sorted[j].PID
	})
	return sorted
}

func sortedProcessChildren(proc *collector.ProcessMetric) []*collector.ProcessMetric {
	children := make([]*collector.ProcessMetric, 0, len(proc.Childs))
	for _, child := range proc.Childs {
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].PID < children[j].PID
	})
	return children
}
