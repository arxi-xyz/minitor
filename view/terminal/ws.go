package terminal

import (
	"context"
	"encoding/json"
	"time"

	"minitor/collector"
	"minitor/transport/socket"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coder/websocket"
)

type snapshotMsg struct {
	CPU          collector.CpuMetric
	RAM          collector.RamMetric
	Disk         collector.DiskMetric
	Network      collector.NetworkMetric
	ProcessCount int
}

type processesMsg struct {
	Processes []*collector.ProcessMetric
}

func listenWS(url string, events chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		go runWS(url, events)
		return nil
	}
}

func waitForWS(events chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-events
	}
}

func runWS(url string, events chan tea.Msg) {
	for {
		_ = runSession(url, events)
		time.Sleep(3 * time.Second)
	}
}

func runSession(url string, events chan tea.Msg) error {
	ctx := context.Background()

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close(websocket.StatusNormalClosure, "") }()

	conn.SetReadLimit(8 << 20)

	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return err
		}

		msg, err := socket.UnmarshalMessage(data)
		if err != nil {
			continue
		}

		switch msg.Type {
		case "snapshot":
			var snap socket.Snapshot
			if err := json.Unmarshal(msg.Data, &snap); err != nil {
				continue
			}

			pushEvent(events, snapshotMsg{
				CPU:          snap.CPU,
				RAM:          snap.RAM,
				Disk:         snap.Disk,
				Network:      snap.Network,
				ProcessCount: snap.ProcessCount,
			})

			if snap.ProcessCount > 0 {
				rows, err := fetchProcesses(ctx, conn, snap.ProcessCount)
				if err != nil {
					return err
				}
				pushEvent(events, processesMsg{Processes: rowsToTree(rows)})
			}

		case "pong":
			continue
		}
	}
}

func fetchProcesses(ctx context.Context, conn *websocket.Conn, total int) ([]socket.ProcessRow, error) {
	rows := make([]socket.ProcessRow, 0, total)

	for offset := 0; offset < total; offset += socket.MaxProcessLimit() {
		reqData, err := json.Marshal(socket.ProcessesRequest{
			Offset: offset,
			Limit:  socket.MaxProcessLimit(),
		})
		if err != nil {
			return nil, err
		}

		req, err := json.Marshal(socket.Message{Type: "processes", Data: reqData})
		if err != nil {
			return nil, err
		}

		if err := conn.Write(ctx, websocket.MessageText, req); err != nil {
			return nil, err
		}

		_, data, err := conn.Read(ctx)
		if err != nil {
			return nil, err
		}

		msg, err := socket.UnmarshalMessage(data)
		if err != nil {
			return nil, err
		}

		var page socket.ProcessesPage
		if err := json.Unmarshal(msg.Data, &page); err != nil {
			return nil, err
		}

		rows = append(rows, page.Items...)
	}

	return rows, nil
}

func rowsToTree(rows []socket.ProcessRow) []*collector.ProcessMetric {
	if len(rows) == 0 {
		return nil
	}

	roots := make([]*collector.ProcessMetric, 0)
	stack := make([]*collector.ProcessMetric, 0)

	for _, row := range rows {
		proc := &collector.ProcessMetric{
			PID:    row.PID,
			PPID:   row.PPID,
			Name:   row.Name,
			Mem:    row.Mem,
			Cpu:    row.Cpu,
			Status: row.Status,
			User:   row.User,
			Childs: make(map[int]*collector.ProcessMetric),
		}

		if row.Depth == 0 {
			roots = append(roots, proc)
			stack = []*collector.ProcessMetric{proc}
			continue
		}

		if row.Depth-1 >= len(stack) {
			continue
		}

		parent := stack[row.Depth-1]
		parent.Childs[proc.PID] = proc
		stack = append(stack[:row.Depth], proc)
	}

	return roots
}

func pushEvent(events chan tea.Msg, msg tea.Msg) {
	select {
	case events <- msg:
	default:
	}
}
