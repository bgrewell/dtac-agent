// +build linux

package network

import (
	"fmt"
	. "github.com/BGrewell/system-api/common"
	"strings"
	"time"
)

var (
	statsCache = make(map[string]*InterfaceStats)
)

type InterfaceStats struct {
	RxBytes   uint64 `json:"rx_bytes"`
	TxBytes   uint64 `json:"tx_bytes"`
	RxPackets uint64 `json:"rx_packets"`
	TxPackets uint64 `json:"tx_packets"`
	RxErrors  uint64 `json:"rx_errors"`
	TxErrors  uint64 `json:"tx_errors"`
	RxDropped uint64 `json:"rx_dropped"`
	TxDropped uint64 `json:"tx_dropped"`
	RxOverrun uint64 `json:"rx_overrun"`
	TxCarrier uint64 `json:"tx_carrier"`
	RxMcast   uint64 `json:"rx_mcast"`
	TxCollsns uint64 `json:"tx_collsns"`
	RxMbps float32 `json:"rx_mbps"`
	TxMbps float32 `json:"tx_mbps"`
	Period float32 `json:"period_sec"`
	recordTime int64
}

func GetInterfaceStats(name string) (stats *InterfaceStats, err error) {
	cmds := []string{fmt.Sprintf("ip -s link show %s", name), "sed -n -e 4p -e 6p"}
	output, err := ExecutePipedCmds(cmds)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 2 {
		return nil, fmt.Errorf("failed to get interface statistics. incorrect output lines")
	}
	rxFields := strings.Fields(lines[0])
	txFields := strings.Fields(lines[1])
	if len(rxFields) != 6 || len(txFields) != 6 {
		return nil, fmt.Errorf("failed to get interface statistics. wrong number of fields")
	}
	var lastRx, lastTx uint64
	var lastRecord int64
	if stats, ok := statsCache[name]; ok {
		lastRx = stats.RxBytes
		lastTx = stats.TxBytes
		lastRecord = stats.recordTime
	}
	stats = &InterfaceStats{
		RxBytes:   ConvertStringToUInt64or0(rxFields[0]),
		TxBytes:   ConvertStringToUInt64or0(txFields[0]),
		RxPackets: ConvertStringToUInt64or0(rxFields[1]),
		TxPackets: ConvertStringToUInt64or0(txFields[1]),
		RxErrors:  ConvertStringToUInt64or0(rxFields[2]),
		TxErrors:  ConvertStringToUInt64or0(txFields[2]),
		RxDropped: ConvertStringToUInt64or0(rxFields[3]),
		TxDropped: ConvertStringToUInt64or0(txFields[3]),
		RxOverrun: ConvertStringToUInt64or0(rxFields[4]),
		TxCarrier: ConvertStringToUInt64or0(txFields[4]),
		RxMcast:   ConvertStringToUInt64or0(rxFields[5]),
		TxCollsns: ConvertStringToUInt64or0(txFields[5]),
	}
	stats.recordTime = time.Now().UnixNano()
	stats.Period = float32(stats.recordTime - lastRecord) / float32(time.Second)
	stats.RxMbps = ConvertToRateMbps(lastRx, stats.RxBytes, lastRecord, stats.recordTime)
	stats.TxMbps = ConvertToRateMbps(lastTx, stats.TxBytes, lastRecord, stats.recordTime)
	statsCache[name] = stats
	return stats, nil
}
