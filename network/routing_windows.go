package network

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// GetRouteTable retrieves the full route table on the system
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
			Destination:    conversions.Inet4_ntoha(entry.dwForwardDest),
			Mask:           conversions.Inet4_ntoha(entry.dwForwardMask),
			NextHop:        conversions.Inet4_ntoha(entry.dwForwardNextHop),
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

// UpdateRoute updates a given route on the system
func UpdateRoute(route RouteTableRow) (err error) {
	return modifyRoute(route, fSetIpForwardEntry)
}

// CreateRoute creates a new route on the system
func CreateRoute(route RouteTableRow) (err error) {
	return modifyRoute(route, fCreateIpForwardEntry)
}

// DeleteRoute removes a route from the system
func DeleteRoute(route RouteTableRow) (err error) {
	return modifyRoute(route, fDeleteIpForwardEntry)
}

// modifyRoute is the core function that handles all changes to routes
func modifyRoute(route RouteTableRow, win32func *windows.LazyProc) (err error) {
	row := MIB_IPFORWARDROW{
		dwForwardDest:      conversions.Inet4_haton(route.Destination),
		dwForwardMask:      conversions.Inet4_haton(route.Mask),
		dwForwardPolicy:    route.Policy,
		dwForwardNextHop:   conversions.Inet4_haton(route.NextHop),
		dwForwardIfIndex:   route.InterfaceIndex,
		dwForwardType:      uint32(route.Type),
		dwForwardProto:     uint32(route.Protocol),
		dwForwardAge:       route.Age,
		dwForwardNextHopAs: route.NextHopAs,
		dwForwardMetric1:   route.Metric1,
		dwForwardMetric2:   route.Metric2,
		dwForwardMetric3:   route.Metric3,
		dwForwardMetric4:   route.Metric4,
		dwForwardMetric5:   route.Metric5,
	}
	ret, _, err := win32func.Call(uintptr(unsafe.Pointer(&row)))
	if ret != windows.NO_ERROR {
		return fmt.Errorf("error: %v", syscall.Errno(ret).Error())
	}
	return nil
}
