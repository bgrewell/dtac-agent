package network

import (
	"encoding/json"
	"fmt"
	"github.com/BGrewell/go-conversions"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type RouteTableRowArgs struct {
	ForceFail bool `json:"force_fail"`
}

// RouteTableRow is the struct for the route table entry
type RouteTableRow struct {
	Destination    string          `json:"destination"`
	Mask           string          `json:"mask"`
	NextHop        string          `json:"next_hop"`
	Policy         uint32          `json:"policy"`
	InterfaceIndex uint32          `json:"interface_index"`
	Type           ForwardType     `json:"type"`
	Protocol       ForwardProtocol `json:"protocol"`
	Age            uint32          `json:"age"`
	NextHopAs      uint32          `json:"next_hop_as"`
	Metric1        uint32          `json:"metric_1"`
	Metric2        uint32          `json:"metric_2"`
	Metric3        uint32          `json:"metric_3"`
	Metric4        uint32          `json:"metric_4"`
	Metric5        uint32          `json:"metric_5"`
}

// String returns the string representation of the route table entry
func (rtr RouteTableRow) String() string {
	return rtr.JSON()
}

// JSON returns the json representation of the route table entry
func (rtr RouteTableRow) JSON() string {
	jout, err := json.Marshal(rtr)
	if err != nil {
		return err.Error()
	}
	return string(jout)
}

// Create creates the route on the system
func (rtr RouteTableRow) Create() error {
	return rtr.modifyRoute(fCreateIpForwardEntry)
}

// Update updates the route on the system
func (rtr RouteTableRow) Update() error {
	return rtr.modifyRoute(fSetIpForwardEntry)
}

// Remove removes the route from the system
func (rtr RouteTableRow) Remove() error {
	return rtr.modifyRoute(fDeleteIpForwardEntry)
}

// Applied returns whether or not the route is applied
func (rtr RouteTableRow) Applied() bool {
	return false
}

// modifyRoute is the core function that handles all changes to routes
func (rtr RouteTableRow) modifyRoute(win32func *windows.LazyProc) (err error) {
	row := MIB_IPFORWARDROW{
		dwForwardDest:      conversions.Inet4_haton(rtr.Destination),
		dwForwardMask:      conversions.Inet4_haton(rtr.Mask),
		dwForwardPolicy:    rtr.Policy,
		dwForwardNextHop:   conversions.Inet4_haton(rtr.NextHop),
		dwForwardIfIndex:   rtr.InterfaceIndex,
		dwForwardType:      uint32(rtr.Type),
		dwForwardProto:     uint32(rtr.Protocol),
		dwForwardAge:       rtr.Age,
		dwForwardNextHopAs: rtr.NextHopAs,
		dwForwardMetric1:   rtr.Metric1,
		dwForwardMetric2:   rtr.Metric2,
		dwForwardMetric3:   rtr.Metric3,
		dwForwardMetric4:   rtr.Metric4,
		dwForwardMetric5:   rtr.Metric5,
	}
	ret, _, err := win32func.Call(uintptr(unsafe.Pointer(&row)))
	if ret != windows.NO_ERROR {
		return fmt.Errorf("error: %v", syscall.Errno(ret).Error())
	}
	return nil
}

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
	return route.Update()
}

// CreateRoute creates a new route on the system
func CreateRoute(route RouteTableRow) (err error) {
	return route.Create()
}

// DeleteRoute removes a route from the system
func DeleteRoute(route RouteTableRow) (err error) {
	return route.Remove()
}
