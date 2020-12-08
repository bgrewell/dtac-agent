// +build linux

package network

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"strings"
	"time"
)

var (
	statsCache = make(map[string]*InterfaceStats)
)

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
		RxBytes:   conversions.ConvertStringToUInt64or0(rxFields[0]),
		TxBytes:   conversions.ConvertStringToUInt64or0(txFields[0]),
		RxPackets: conversions.ConvertStringToUInt64or0(rxFields[1]),
		TxPackets: conversions.ConvertStringToUInt64or0(txFields[1]),
		RxErrors:  conversions.ConvertStringToUInt64or0(rxFields[2]),
		TxErrors:  conversions.ConvertStringToUInt64or0(txFields[2]),
		RxDropped: conversions.ConvertStringToUInt64or0(rxFields[3]),
		TxDropped: conversions.ConvertStringToUInt64or0(txFields[3]),
		RxOverrun: conversions.ConvertStringToUInt64or0(rxFields[4]),
		TxCarrier: conversions.ConvertStringToUInt64or0(txFields[4]),
		RxMcast:   conversions.ConvertStringToUInt64or0(rxFields[5]),
		TxCollsns: conversions.ConvertStringToUInt64or0(txFields[5]),
	}
	stats.recordTime = time.Now().UnixNano()
	stats.Period = float32(stats.recordTime-lastRecord) / float32(time.Second)
	stats.RxMbps = conversions.ConvertToRateMbps(lastRx, stats.RxBytes, lastRecord, stats.recordTime)
	stats.TxMbps = conversions.ConvertToRateMbps(lastTx, stats.TxBytes, lastRecord, stats.recordTime)
	statsCache[name] = stats
	return stats, nil
}
