package network

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"net"
)

type RouteTableRow struct {
	LinkIndex  int                    `json:"link_index"`
	ILinkIndex int                    `json:"i_link_index"`
	Scope      netlink.Scope          `json:"scope"`
	Dst        string                 `json:"dst"`
	DstMask    string                 `json:"dst_mask"`
	Src        net.IP                 `json:"src"`
	Gw         net.IP                 `json:"gw"`
	MultiPath  []*netlink.NexthopInfo `json:"multi_path"`
	Protocol   int                    `json:"protocol"`
	Priority   int                    `json:"priority"`
	Table      int                    `json:"table"`
	Type       int                    `json:"type"`
	Tos        int                    `json:"tos"`
	Flags      int                    `json:"flags"`
	MPLSDst    *int                   `json:"mpls_dst"`
	NewDst     string                 `json:"new_dst"`
	Encap      netlink.Encap          `json:"encap"`
	MTU        int                    `json:"mtu"`
	AdvMSS     int                    `json:"adv_mss"`
	Hoplimit   int                    `json:"hoplimit"`
}

// GetRouteTable retrieves the full route table on the system
func GetRouteTable() (routes []RouteTableRow, err error) {
	return nil, fmt.Errorf("this method has not been implemented for this OS")
}

// UpdateRoute updates a given route on the system
func UpdateRoute(route RouteTableRow) (err error) {
	return fmt.Errorf("this method has not been implemented for this OS")
}

// CreateRoute creates a new route on the system
func CreateRoute(route RouteTableRow) (err error) {
	return fmt.Errorf("this method has not been implemented for this OS")
}

// DeleteRoute removes a route from the system
func DeleteRoute(route RouteTableRow) (err error) {
	return fmt.Errorf("this method has not been implemented for this OS")
}
