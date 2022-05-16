package common

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type BasicInfo struct {
	Host      *host.InfoStat            `json:"host"`
	CPU       []cpu.InfoStat            `json:"cpu"`
	Memory    *mem.VirtualMemoryStat    `json:"memory"`
	Disk      *disk.UsageStat           `json:"disk"`
	Network   []net.InterfaceStat       `json:"network"`
	Docker    []docker.CgroupDockerStat `json:"docker"`
	Processes []*ProcessInfo            `json:"processes"`
	Routes    EndpointsInfoParser       `json:"http_endpoint_info"`
}

type EndpointInfo struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"-"`
}

type EndpointsInfoParser struct {
	Routes []*EndpointInfo `json:"endpoints"`
}

type ProcessInfo struct {
	Pid        int32                   `json:"pid"`
	Ppid       int32                   `json:"ppid"`
	Name       string                  `json:"name"`
	Cmdline    string                  `json:"cmdline"`
	CreateTime int64                   `json:"created_time"`
	Exe        string                  `json:"exe"`
	IoCounters *process.IOCountersStat `json:"io_counters"`
	Nice       int32                   `json:"nice"`
	NumThreads int32                   `json:"num_threads"`
	MemoryInfo *process.MemoryInfoStat `json:"memory_info"`
	Username   string                  `json:"username"`
}

func (pi *ProcessInfo) Update(p *process.Process) (err error) {
	pi.Pid = p.Pid
	pi.Ppid, _ = p.Ppid()
	pi.Name, _ = p.Name()
	pi.Cmdline, _ = p.Cmdline()
	pi.CreateTime, _ = p.CreateTime()
	pi.Exe, _ = p.Exe()
	pi.IoCounters, _ = p.IOCounters()
	pi.Nice, _ = p.Nice()
	pi.NumThreads, _ = p.NumThreads()
	pi.MemoryInfo, _ = p.MemoryInfo()
	pi.Username, _ = p.Username()
	return nil
}

func (bi *BasicInfo) Update(routes gin.RoutesInfo) {
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

	// Update Disk Info
	bi.Disk, err = disk.Usage("/")
	if err != nil {
		bi.Disk = nil
	}

	// Update Docker Info
	bi.Docker, err = docker.GetDockerStat()
	if err != nil {
		bi.Docker = nil
	}

	// Update Network Info
	bi.Network, err = net.Interfaces()
	if err != nil {
		bi.Network = nil
	}

	// Update Process Info
	bi.Processes = make([]*ProcessInfo, 0)
	processes, err := process.Processes()
	if err != nil {
		bi.Processes = nil
	} else {
		for _, p := range processes {
			pi := ProcessInfo{}
			pi.Update(p)
			bi.Processes = append(bi.Processes, &pi)
		}
	}

	// UpdateRoutes
	bi.UpdateRoutes(routes)

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
