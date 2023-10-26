package hardware

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
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
	name       string // Subsystem name
	nic        NicInfo
	cpu        CPUInfo
	mem        MemoryInfo
	disk       DiskInfo
	endpoints  []endpoint.Endpoint
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
	s.endpoints = []endpoint.Endpoint{
		{fmt.Sprintf("%s/cpu", base), endpoint.ActionRead, s.cpuInfoHandler, secure, nil, nil},
		{fmt.Sprintf("%s/cpu/usage", base), endpoint.ActionRead, s.cpuUsageHandler, secure, CpuUsageArgs{}, nil},

		{fmt.Sprintf("%s/memory", base), endpoint.ActionRead, s.memInfoHandler, secure, nil, nil},

		{fmt.Sprintf("%s/disk", base), endpoint.ActionRead, s.diskRootHandler, secure, nil, nil},
		{fmt.Sprintf("%s/disk/partitions", base), endpoint.ActionRead, s.diskPartitionHandler, secure, nil, nil},
		{fmt.Sprintf("%s/disk/disks", base), endpoint.ActionRead, s.diskPhysicalDisksHandler, secure, nil, nil},
		{fmt.Sprintf("%s/disk/usage", base), endpoint.ActionRead, s.diskUsageHandler, secure, DiskUsageArgs{}, nil},

		{fmt.Sprintf("%s/network", base), endpoint.ActionRead, s.nicRootHandler, secure, nil, nil},
		// TODO: Support query param 'name' to get specific interface
		{fmt.Sprintf("%s/network/interfaces", base), endpoint.ActionRead, s.nicRootHandler, secure, NicInfoArgs{}, nil},
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
func (s *Subsystem) Endpoints() []endpoint.Endpoint {
	return s.endpoints
}
