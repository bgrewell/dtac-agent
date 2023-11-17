package network

import (
	"encoding/json"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
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

func (s *Subsystem) arpTableHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		arpData, err := GetArpTable()
		if err != nil {
			return nil, err
		}
		return json.Marshal(arpData)
	}, "returns ARP table information from the system")
}
