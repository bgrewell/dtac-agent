package hardware

import (
	"encoding/json"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
)

// NicInfoArgs is a struct to assist with validating the input arguments
type NicInfoArgs struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
}

// NicInfo is the interface for the nic subsystem
type NicInfo interface {
	Update()
	Info() []net.InterfaceStat
}

// LiveNicInfo is the struct for the nic subsystem
type LiveNicInfo struct {
	Logger         *zap.Logger // All subsystems have a pointer to the logger
	InterfaceStats []net.InterfaceStat
}

// Update updates the nic subsystem
func (ni *LiveNicInfo) Update() {
	n, err := net.Interfaces()
	if err != nil {
		ni.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	ni.InterfaceStats = n
}

// Info returns the nic subsystem info
func (ni *LiveNicInfo) Info() []net.InterfaceStat {
	return ni.InterfaceStats
}

// rootHandler handles requests for the root path for this subsystem
func (s *Subsystem) nicRootHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		name := ""
		if v, ok := in.Parameters["name"]; ok {
			name = v[0]
		}

		s.nic.Update()
		if name == "" {
			return json.Marshal(s.nic.Info())
		}

		for _, info := range s.nic.Info() {
			if info.Name == name {
				return json.Marshal(info)
			}
		}

		return nil, fmt.Errorf("no interface found by name: %s", name)
	}, "network interface information")
}
