// +build darwin

package network

import (
	"fmt"
)

func GetInterfaceStats(name string) (stats *InterfaceStats, err error) {
	return nil, fmt.Errorf("this function has not been implemented on macOS yet")
}
