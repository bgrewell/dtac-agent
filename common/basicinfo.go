package common

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type BasicInfo struct {
	Host    *host.InfoStat         `json:"host"`
	CPU     []cpu.InfoStat         `json:"cpu"`
	Memory  *mem.VirtualMemoryStat `json:"memory"`
	Network []net.InterfaceStat    `json:"network"`
	Routes  EndpointsInfoParser    `json:"http_endpoint_info"`
}

type EndpointInfo struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"-"`
}

type EndpointsInfoParser struct {
	Routes []*EndpointInfo `json:"endpoints"`
}

func (bi *BasicInfo) Update() {
	// Update Host
	var err error
	bi.Host, err = host.Info()
	if err != nil {
		bi.Host = nil
	}

	// Update CPU Info
	bi.CPU, err = cpu.Info()
	if err != nil {
		bi.CPU = nil
	}

	// Update Mem Info
	bi.Memory, err = mem.VirtualMemory()
	if err != nil {
		bi.Memory = nil
	}

	// Update Network Info
	bi.Network, err = net.Interfaces()
}

func (bi *BasicInfo) UpdateRoutes(routes gin.RoutesInfo) {
	bi.Routes = EndpointsInfoParser{Routes: make([]*EndpointInfo, len(routes))}
	for idx, route := range routes {
		bi.Routes.Routes[idx] = &EndpointInfo{
			Method:  route.Method,
			Path:    route.Path,
			Handler: route.Handler,
		}
	}
}
