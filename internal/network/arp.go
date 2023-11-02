package network

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
)

// ArpEntry is the struct for the arp entry
type ArpEntry struct {
	IPAddress string `json:"ip_address"`
	HWType    string `json:"hw_type"`
	Flags     string `json:"flags"`
	HWAddress string `json:"hw_address"`
	Mask      string `json:"mask"`
	Iface     string `json:"device"`
}

func (s *Subsystem) arpTableHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		arpData, err := GetArpTable()
		if err != nil {
			return nil, err
		}
		return arpData, nil
	}, "returns ARP table information from the system")
}
