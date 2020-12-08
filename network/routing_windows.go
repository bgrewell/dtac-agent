package network

import (
	"fmt"
	"github.com/BGrewell/system-api/common"
	"unsafe"
)

func GetRouteTable() (routes []RouteTableRow, err error) {
	ft := MIB_IPFORWARDTABLE{}
	dwSize := uint32(0)
	ret, _, err := fGetIpForwardTable.Call(uintptr(unsafe.Pointer(&ft)), uintptr(unsafe.Pointer(&dwSize)), 0)
	if ret != 0 && ret != 122 {
		fmt.Errorf("failed to get size of route table: %v", err)
	}
	//entries := dwSize / MIB_IPFORWARDROW{}.Size()
	ft.table = [10000]MIB_IPFORWARDROW{}
	ret, _, err = fGetIpForwardTable.Call(uintptr(unsafe.Pointer(&ft)), uintptr(unsafe.Pointer(&dwSize)), 0)
	if ret != 0 {
		fmt.Errorf("failed to get route table: %v", err)
	}
	rows := make([]RouteTableRow, ft.dwNumEntries)
	for i := 0; i < int(ft.dwNumEntries); i++ {
		entry := ft.table[i]
		row := RouteTableRow{
			Destination:    common.Inet_ntoha(entry.dwForwardDest),
			Mask:           common.Inet_ntoha(entry.dwForwardMask),
			NextHop:        common.Inet_ntoha(entry.dwForwardNextHop),
			Policy:         entry.dwForwardPolicy,
			InterfaceIndex: entry.dwForwardIfIndex,
			Type:           ForwardType(entry.dwForwardType),
			Protocol:       ForwardProtocol(entry.dwForwardProto),
			Age:            entry.dwForwardAge,
			NextHopAs:      entry.dwForwardNextHopAs,
			Metric1:        entry.dwForwardMetric1,
			Metric2:        entry.dwForwardMetric2,
			Metric3:        entry.dwForwardMetric3,
			Metric4:        entry.dwForwardMetric4,
			Metric5:        entry.dwForwardMetric5,
		}
		rows[i] = row
	}

	return rows, nil
}
