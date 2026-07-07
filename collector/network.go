package collector

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	gopsutilnet "github.com/shirou/gopsutil/v4/net"
)

type NetworkMetric struct {
	InterfaceName string
	IP            string

	RxBytes uint64
	TxBytes uint64

	PacketsRx uint64
	PacketsTx uint64

	Errors uint64
	Drops  uint64
}

func GetNetworkInfo() (NetworkMetric, error) {
	stats, err := gopsutilnet.IOCounters(true)
	if err != nil {
		return NetworkMetric{}, err
	}

	statsByName := make(map[string]gopsutilnet.IOCountersStat, len(stats))
	for _, stat := range stats {
		statsByName[stat.Name] = stat
	}

	var best gopsutilnet.IOCountersStat
	found := false

	if defaultIface := getDefaultRouteInterface(); defaultIface != "" {
		if stat, ok := statsByName[defaultIface]; ok && !isVirtualInterface(stat.Name) {
			best = stat
			found = true
		}
	}

	if !found {
		for _, stat := range stats {
			if isVirtualInterface(stat.Name) {
				continue
			}

			if !found || (stat.BytesRecv+stat.BytesSent) > (best.BytesRecv+best.BytesSent) {
				best = stat
				found = true
			}
		}
	}

	if !found {
		return NetworkMetric{}, nil
	}

	return NetworkMetric{
		InterfaceName: best.Name,
		IP:            getInterfaceIP(best.Name),
		RxBytes:       best.BytesRecv,
		TxBytes:       best.BytesSent,
		PacketsRx:     best.PacketsRecv,
		PacketsTx:     best.PacketsSent,
		Errors:        best.Errin + best.Errout,
		Drops:         best.Dropin + best.Dropout,
	}, nil
}

func getDefaultRouteInterface() string {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		dest, err := strconv.ParseInt(fields[1], 16, 64)
		if err != nil || dest != 0 {
			continue
		}

		return fields[0]
	}

	return ""
}

func isVirtualInterface(name string) bool {
	prefixes := []string{"lo", "docker", "br-", "veth", "virbr", "tun", "tap"}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func getInterfaceIP(interfaceName string) string {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		ip := ipNet.IP

		if ip.IsLoopback() {
			continue
		}

		if ip.To4() != nil {
			return ip.String()
		}
	}

	return ""
}

func WorkerNetwork(ch chan<- NetworkMetric) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metric, err := GetNetworkInfo()
		if err != nil {
			continue
		}

		ch <- metric
	}
}
