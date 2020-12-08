package network

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"
)

type InterfaceStats struct {
	RxBytes    uint64  `json:"rx_bytes"`
	TxBytes    uint64  `json:"tx_bytes"`
	RxPackets  uint64  `json:"rx_packets"`
	TxPackets  uint64  `json:"tx_packets"`
	RxErrors   uint64  `json:"rx_errors"`
	TxErrors   uint64  `json:"tx_errors"`
	RxDropped  uint64  `json:"rx_dropped"`
	TxDropped  uint64  `json:"tx_dropped"`
	RxOverrun  uint64  `json:"rx_overrun"`
	TxCarrier  uint64  `json:"tx_carrier"`
	RxMcast    uint64  `json:"rx_mcast"`
	TxCollsns  uint64  `json:"tx_collsns"`
	RxMbps     float32 `json:"rx_mbps"`
	TxMbps     float32 `json:"tx_mbps"`
	Period     float32 `json:"period_sec"`
	recordTime int64
}

// Addr is an override to control how net.Addr is marshaled to json
type Addr struct {
	IP      string `json:"ip"`
	Network string `json:"network"`
}

// Interface is an override to control how the net.Interface is marshalled to json
type Interface struct {
	Index              int             `json:"index"`         // positive integer that starts at one, zero is never used
	MTU                int             `json:"mtu"`           // maximum transmission unit
	Name               string          `json:"name"`          // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr       string          `json:"hardware_addr"` // IEEE MAC-48, EUI-48 and EUI-64 form
	Flags              string          `json:"flags"`         // e.g., FlagUp, FlagLoopback, FlagMulticast
	FlagsInt           int             `json:"flags_int"`
	Addresses          []*Addr         `json:"addresses"`
	MulticastAddresses []*Addr         `json:"multicast_addresses"`
	Statistics         *InterfaceStats `json:"statistics"`
}

func GetInterfaces() (ifaces []*Interface, err error) {
	interfaces, err := net.Interfaces()
	ifaces = make([]*Interface, len(interfaces))
	if err != nil {
		return nil, err
	}
	for idx, iface := range interfaces {

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		addresses := make([]*Addr, len(addrs))
		for idx, addr := range addrs {
			addresses[idx] = &Addr{
				IP:      addr.String(),
				Network: addr.Network(),
			}
		}

		maddrs, err := iface.MulticastAddrs()
		if err != nil {
			return nil, err
		}
		maddresses := make([]*Addr, len(maddrs))
		for idx, addr := range maddrs {
			maddresses[idx] = &Addr{
				IP:      addr.String(),
				Network: addr.Network(),
			}
		}
		stats, err := GetInterfaceStats(iface.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"name":  iface.Name,
				"stats": stats,
				"err":   err,
			}).Debug("failed to get interface statistics")
		}
		i := &Interface{
			Index:              iface.Index,
			MTU:                iface.MTU,
			Name:               iface.Name,
			HardwareAddr:       iface.HardwareAddr.String(),
			Flags:              iface.Flags.String(),
			FlagsInt:           int(iface.Flags),
			Addresses:          addresses,
			MulticastAddresses: maddresses,
			Statistics:         stats,
		}
		ifaces[idx] = i
	}
	return ifaces, nil
}

func GetInterfaceNames() (ifaces []string, err error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(interfaces))
	for idx, iface := range interfaces {
		names[idx] = iface.Name
	}
	return names, nil
}

func GetInterfaceByName(name string) (iface *Interface, err error) {
	ifaces, err := GetInterfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Name == name {
			return iface, nil
		}
	}
	return nil, fmt.Errorf("no interface with that name was found")
}

func GetInterfaceByIdx(idStr string) (iface *Interface, err error) {
	ifaces, err := GetInterfaces()
	if err != nil {
		return nil, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Index == id {
			return iface, nil
		}
	}
	return nil, fmt.Errorf("no interface with that id was found")
}
