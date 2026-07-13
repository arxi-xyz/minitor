package socket

import (
	"log"
	"sync"

	"minitor/collector"
)

type Monitor struct {
	hub  *Hub
	snap Snapshot

	processMu sync.RWMutex
	flat      []ProcessRow
}

func NewMonitor(hub *Hub) *Monitor {
	return &Monitor{hub: hub}
}

func (m *Monitor) Run() {
	cpuCh := make(chan collector.CpuMetric)
	ramCh := make(chan collector.RamMetric)
	diskCh := make(chan collector.DiskMetric)
	networkCh := make(chan collector.NetworkMetric)
	processCh := make(chan []*collector.ProcessMetric)

	go collector.WorkerCpu(cpuCh)
	go collector.WorkerRam(ramCh)
	go collector.WorkerDisk(diskCh)
	go collector.WorkerNetwork(networkCh)
	go collector.ProcessWorker(processCh)

	for {
		select {
		case v := <-cpuCh:
			m.snap.CPU = v
			m.broadcast()
		case v := <-ramCh:
			m.snap.RAM = v
			m.broadcast()
		case v := <-diskCh:
			m.snap.Disk = v
			m.broadcast()
		case v := <-networkCh:
			m.snap.Network = v
			m.broadcast()
		case v := <-processCh:
			m.setProcesses(v)
			m.broadcast()
		}
	}
}

func (m *Monitor) setProcesses(processes []*collector.ProcessMetric) {
	flat := flattenProcesses(processes)

	m.processMu.Lock()
	m.flat = flat
	m.snap.ProcessCount = len(flat)
	m.processMu.Unlock()
}

func (m *Monitor) ProcessesPage(offset, limit int) ProcessesPage {
	m.processMu.RLock()
	defer m.processMu.RUnlock()

	total := len(m.flat)
	if limit <= 0 {
		limit = defaultProcessLimit
	}
	if limit > maxProcessLimit {
		limit = maxProcessLimit
	}
	if offset < 0 {
		offset = 0
	}
	if offset > total {
		offset = total
	}

	end := offset + limit
	if end > total {
		end = total
	}

	items := make([]ProcessRow, end-offset)
	copy(items, m.flat[offset:end])

	return ProcessesPage{
		Offset: offset,
		Limit:  limit,
		Total:  total,
		Items:  items,
	}
}

func (m *Monitor) broadcast() {
	msg, err := m.snap.Message()
	if err != nil {
		log.Println(err)
		return
	}

	m.hub.Broadcast(msg)
}
