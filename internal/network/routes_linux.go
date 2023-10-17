package network

import (
	"encoding/json"
	"fmt"
	"github.com/BGrewell/go-conversions"
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

func (rtr RouteTableRow) String() string {
	return rtr.JSON()
}

func (rtr RouteTableRow) JSON() string {
	jout, err := json.Marshal(rtr)
	if err != nil {
		return err.Error()
	}
	return string(jout)
}

func (rtr RouteTableRow) Create() error {
	return rtr.modifyRoute(Route_Create)
}

func (rtr RouteTableRow) Update() error {
	return rtr.modifyRoute(Route_Update)
}

func (rtr RouteTableRow) Remove() error {
	return rtr.modifyRoute(Route_Delete)
}

func (rtr RouteTableRow) Applied() bool {
	return false
}

func (rtr RouteTableRow) modifyRoute(action RouteAction) (err error) {
	var Dst *net.IPNet = nil
	if rtr.Dst != "" && rtr.DstMask != "" {
		mask, err := conversions.Ipv4MaskBytes(rtr.DstMask)
		if err != nil {
			return err
		}
		Dst = &net.IPNet{
			IP:   net.ParseIP(rtr.Dst),
			Mask: mask,
		}
	}
	internalRoute := netlink.Route{
		LinkIndex:  rtr.LinkIndex,
		ILinkIndex: rtr.ILinkIndex,
		Scope:      rtr.Scope,
		Dst:        Dst,
		Src:        rtr.Src,
		Gw:         rtr.Gw,
		MultiPath:  rtr.MultiPath,
		Protocol:   rtr.Protocol,
		Priority:   rtr.Priority,
		Table:      rtr.Table,
		Type:       rtr.Type,
		Tos:        rtr.Tos,
		Flags:      rtr.Flags,
		MPLSDst:    rtr.MPLSDst,
		NewDst:     nil, // todo: not supported yet, the linux implementation in general needs to be refactored
		Encap:      rtr.Encap,
		MTU:        rtr.MTU,
		AdvMSS:     rtr.AdvMSS,
		Hoplimit:   rtr.Hoplimit,
	}
	switch action {
	case Route_Create:
		return netlink.RouteAdd(&internalRoute)
	case Route_Update:
		return netlink.RouteReplace(&internalRoute)
	case Route_Delete:
		return netlink.RouteDel(&internalRoute)
	default:
		return fmt.Errorf("unknown route action")
	}
}

type RouteAction int

const (
	Route_Create RouteAction = 1
	Route_Update RouteAction = 2
	Route_Delete RouteAction = 3
)

// GetRouteTable retrieves the full route table on the system
func GetRouteTable() (routes []RouteTableRow, err error) {
	internalRoutes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	routes = make([]RouteTableRow, len(internalRoutes))
	for idx, route := range internalRoutes {
		NewDst := ""
		if route.NewDst != nil {
			NewDst = route.NewDst.String()
		}
		Dst := ""
		DstMask := ""
		if route.Dst != nil {
			mask, err := conversions.Ipv4MaskString(route.Dst.Mask)
			if err != nil {
				return nil, err
			}
			Dst = route.Dst.IP.String()
			DstMask = mask
		}
		rtr := RouteTableRow{
			LinkIndex:  route.LinkIndex,
			ILinkIndex: route.ILinkIndex,
			Scope:      route.Scope,
			Dst:        Dst,
			DstMask:    DstMask,
			Src:        route.Src,
			Gw:         route.Gw,
			MultiPath:  route.MultiPath,
			Protocol:   route.Protocol,
			Priority:   route.Priority,
			Table:      route.Table,
			Type:       route.Type,
			Tos:        route.Tos,
			Flags:      route.Flags,
			MPLSDst:    route.MPLSDst,
			NewDst:     NewDst,
			Encap:      route.Encap,
			MTU:        route.MTU,
			AdvMSS:     route.AdvMSS,
			Hoplimit:   route.Hoplimit,
		}
		routes[idx] = rtr
	}
	return routes, nil
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
