// +build windows

package network

import (
	"fmt"
	"strings"
	"github.com/BGrewell/go-execute"
)

var (
	statsCache = make(map[string]*InterfaceStats)
)

func GetInterfaceStats(name string) (stats *InterfaceStats, err error) {
	cmd := fmt.Sprintf("Get-NetAdapter -Name %s | Get-NetAdapterStatistics | Format-List -Property \"*\"", name)
	output, stderr, err := execute.ExecutePowershell(cmd)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		fmt.Println(line)
	}
	return nil, fmt.Errorf("this function has not been implemented on Windows yet")
}
