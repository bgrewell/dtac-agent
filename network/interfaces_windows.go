// +build windows

package network

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"strconv"
	"strings"
	"github.com/BGrewell/go-execute"
	"time"
)

var (
	statsCache = make(map[string]*InterfaceStats)
)

func GetInterfaceStats(name string) (stats *InterfaceStats, err error) {
	// Execute the powershell command to get the stats
	cmd := fmt.Sprintf("Get-NetAdapter -Name %s | Get-NetAdapterStatistics | Format-List -Property \"*\"", name)
	output, stderr, err := execute.ExecutePowershell(cmd)
	if err != nil {
		fmt.Println(stderr)
		return nil, err
	}

	// Parse the output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	elements := make(map[string]string)
	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			elements[parts[0]] = parts[2]
		}
	}

	// Convert the fields we are interested in
	rxBytes, _ := strconv.Atoi(elements["ReceivedBytes"])
	txBytes, _ := strconv.Atoi(elements["SentBytes"])
	rxPacketsUni, _ := strconv.Atoi(elements["ReceivedUnicastPackets"])
	txPacketsUni, _ := strconv.Atoi(elements["SentUnicastPackets"])
	rxPacketsMulti, _ := strconv.Atoi(elements["ReceivedMulticastPackets"])
	txPacketsMulti, _ := strconv.Atoi(elements["SentMulticastPackets"])
	rxPacketsBroad, _ := strconv.Atoi(elements["ReceivedBroadcastPackets"])
	txPacketsBroad, _ := strconv.Atoi(elements["SentBroadcastPackets"])
	rxErrors, _ := strconv.Atoi(elements["ReceivedPacketErrors"])
	txErrors, _ := 0, 0 // unsupported on windows
	rxDropped, _ := strconv.Atoi(elements["ReceivedDiscardedPackets"])
	txDropped, _ := 0, 0 // unsupported on windows

	// Populate the last values used for calculated fields if they exist in the cache
	var lastRx, lastTx uint64
	var lastRecord int64
	if stats, ok := statsCache[name]; ok {
		lastRx = stats.RxBytes
		lastTx = stats.TxBytes
		lastRecord = stats.recordTime
	}

	// Fill out fields from above values
	stats = &InterfaceStats{
		RxBytes:    uint64(rxBytes),
		TxBytes:    uint64(txBytes),
		RxPackets:  uint64(rxPacketsUni + rxPacketsMulti + rxPacketsBroad),
		TxPackets:  uint64(txPacketsUni + txPacketsMulti + txPacketsBroad),
		RxErrors:   uint64(rxErrors),
		TxErrors:   uint64(txErrors),
		RxDropped:  uint64(rxDropped),
		TxDropped:  uint64(txDropped),
		RxOverrun:  0,
		TxCarrier:  0,
		RxMcast:    0,
		TxCollsns:  0,
	}
	// Fill calculated fields
	stats.recordTime = time.Now().UnixNano()
	stats.Period = float32(stats.recordTime-lastRecord) / float32(time.Second)
	stats.RxMbps = conversions.ConvertToRateMbps(lastRx, stats.RxBytes, lastRecord, stats.recordTime)
	stats.TxMbps = conversions.ConvertToRateMbps(lastTx, stats.TxBytes, lastRecord, stats.recordTime)

	// Update the stats cache so it can be used in the next calculations
	statsCache[name] = stats
	return stats, nil
}
