package hardware

import (
	"fmt"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"github.com/bgrewell/dtac-agent/internal/controller"
	"github.com/bgrewell/dtac-agent/internal/interfaces"
	"go.uber.org/zap"
)

// NewSubsystem creates a new instance of the Subsystem and if that subsystem is enabled it calls
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "hardware"
	logger := c.Logger.With(zap.String("module", name))
	hw := Subsystem{
		Controller: c,
		Logger:     logger,
		enabled:    c.Config.Subsystems.Diag,
		name:       name,
		nic:        &LiveNicInfo{Logger: logger},
		cpu:        &LiveCPUInfo{Logger: logger},
		mem:        &LiveMemoryInfo{Logger: logger},
		disk:       &LiveDiskInfo{Logger: logger},
	}
	hw.register()
	return &hw
}

// Subsystem is the subsystem that contains routes related to internal dtac diagnostics
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	nic        NicInfo
	cpu        CPUInfo
	mem        MemoryInfo
	disk       DiskInfo
	endpoints  []*endpoint.Endpoint
}

// register registers the routes that this module handles
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	// Create a group for this subsystem
	base := s.name

	// Endpoints
	secure := s.Controller.Config.Auth.DefaultSecure
	authz := endpoint.AuthGroupUser.String()
	s.endpoints = []*endpoint.Endpoint{
		endpoint.NewEndpoint(fmt.Sprintf("%s/cpu", base), endpoint.ActionRead, "cpu information", s.cpuInfoHandler, secure, authz, endpoint.WithOutput([]cpu.InfoStat{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/cpu/usage", base), endpoint.ActionRead, "cpu usage information", s.cpuUsageHandler, secure, authz, endpoint.WithParameters(CPUUsageArgs{}), endpoint.WithOutput(CPUUsageOutput{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/memory", base), endpoint.ActionRead, "memory information", s.memInfoHandler, secure, authz, endpoint.WithOutput(&mem.VirtualMemoryStat{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/disk", base), endpoint.ActionRead, "disk information", s.diskRootHandler, secure, authz, endpoint.WithOutput(&DiskReport{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/disk/partitions", base), endpoint.ActionRead, "disk partition information", s.diskPartitionHandler, secure, authz, endpoint.WithOutput([]disk.PartitionStat{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/disk/disks", base), endpoint.ActionRead, "list of physical disks", s.diskPhysicalDisksHandler, secure, authz, endpoint.WithOutput([]*DiskDetails{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/disk/usage", base), endpoint.ActionRead, "disk usage", s.diskUsageHandler, secure, authz, endpoint.WithParameters(DiskUsageArgs{}), endpoint.WithOutput([]*disk.UsageStat{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/network", base), endpoint.ActionRead, "network information", s.nicRootHandler, secure, authz, endpoint.WithOutput([]net.InterfaceStat{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/network/interfaces", base), endpoint.ActionRead, "network interface information", s.nicRootHandler, secure, authz, endpoint.WithParameters(NicInfoArgs{}), endpoint.WithOutput([]net.InterfaceStat{})),
	}

}

// Enabled returns true if this module is enabled otherwise it returns false
func (s *Subsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the subsystem
func (s *Subsystem) Name() string {
	return s.name
}

// Endpoints returns an array of endpoints that this Subsystem handles
func (s *Subsystem) Endpoints() []*endpoint.Endpoint {
	return s.endpoints
}
