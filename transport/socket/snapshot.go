package socket

import (
	"encoding/json"

	"minitor/collector"
)

type Snapshot struct {
	CPU          collector.CpuMetric     `json:"cpu"`
	RAM          collector.RamMetric     `json:"ram"`
	Disk         collector.DiskMetric    `json:"disk"`
	Network      collector.NetworkMetric `json:"network"`
	ProcessCount int                     `json:"process_count"`
}

func (s Snapshot) Message() (Message, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return Message{}, err
	}

	return Message{Type: "snapshot", Data: data}, nil
}
