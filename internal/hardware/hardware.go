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
		{Path: fmt.Sprintf("%s/cpu", base), Action: endpoint.ActionRead, Function: s.cpuInfoHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/cpu/usage", base), Action: endpoint.ActionRead, Function: s.cpuUsageHandler, UsesAuth: secure, ExpectedArgs: CpuUsageArgs{}, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/memory", base), Action: endpoint.ActionRead, Function: s.memInfoHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/disk", base), Action: endpoint.ActionRead, Function: s.diskRootHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/disk/partitions", base), Action: endpoint.ActionRead, Function: s.diskPartitionHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/disk/disks", base), Action: endpoint.ActionRead, Function: s.diskPhysicalDisksHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/disk/usage", base), Action: endpoint.ActionRead, Function: s.diskUsageHandler, UsesAuth: secure, ExpectedArgs: DiskUsageArgs{}, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/network", base), Action: endpoint.ActionRead, Function: s.nicRootHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/network/interfaces", base), Action: endpoint.ActionRead, Function: s.nicRootHandler, UsesAuth: secure, ExpectedArgs: NicInfoArgs{}, ExpectedBody: nil},
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
