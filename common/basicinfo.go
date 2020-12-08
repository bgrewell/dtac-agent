package common

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type BasicInfo struct {
	Host *host.InfoStat `json:"host"`
	CPU []cpu.InfoStat `json:"cpu"`
	Memory *mem.VirtualMemoryStat `json:"memory"`
	Network []net.InterfaceStat `json:"network"`
	Routes RoutesInfoParser `json:"route_info"`
}

type RouteInfoParser struct {
	Method string `json:"method"`
	Path string `json:"path"`
	Handler string `json:"-"`
}

type RoutesInfoParser struct {
	Routes []*RouteInfoParser
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
	bi.Routes = RoutesInfoParser{Routes: make([]*RouteInfoParser, len(routes))}
	for idx, route := range routes {
		bi.Routes.Routes[idx] = &RouteInfoParser{
			Method:  route.Method,
			Path:    route.Path,
			Handler: route.Handler,
		}
	}
}
