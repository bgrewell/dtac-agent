package network

import (
	"fmt"
	"unsafe"
)

func GetIpStatistics() (err error) {
	if fGetIpStatistics.Find() != nil {
		fmt.Errorf("GetIpStatistics not found")
	}
	IpStats := MIB_IP_STATS_LH{}
	ret, _, err := fGetIpStatistics.Call(uintptr(unsafe.Pointer(&IpStats)))
	if ret != 0 {
		fmt.Println("Failed to get IP Stats")
	}
	return nil
}
