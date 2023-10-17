package hardware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/register"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
)

// NewSubsystem creates a new instance of the HardwareSubsystem and if that subsystem is enabled it calls
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "hardware"
	logger := c.Logger.With(zap.String("module", name))
	hw := HardwareSubsystem{
		Controller: c,
		Logger:     logger,
		enabled:    c.Config.Subsystems.Diag,
		name:       name,
		nic:        &LiveNicInfo{Logger: logger},
		cpu:        &LiveCpuInfo{Logger: logger},
		mem:        &LiveMemoryInfo{Logger: logger},
		disk:       &LiveDiskInfo{Logger: logger},
	}
	return &hw
}

// HardwareSubsystem is the subsystem that contains routes related to internal dtac diagnostics
type HardwareSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string // Subsystem name
	nic        NicInfo
	cpu        CpuInfo
	mem        MemoryInfo
	disk       DiskInfo
}

// Register registers the routes that this module handles
func (s *HardwareSubsystem) Register() error {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return nil
	}
	// Create a group for this subsystem
	base := s.Controller.Router.Group(s.name)

	// Routes
	secure := s.Controller.Config.Auth.DefaultSecure
	routes := []types.RouteInfo{
		// CPU Routes
		{Group: base, HttpMethod: http.MethodGet, Path: "/cpu", Handler: s.cpuInfoHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/cpu/usage", Handler: s.cpuUsageHandler, Protected: secure},
		// Memory Routes
		{Group: base, HttpMethod: http.MethodGet, Path: "/memory", Handler: s.memInfoHandler, Protected: secure},
		// Disk Routes
		{Group: base, HttpMethod: http.MethodGet, Path: "/disk", Handler: s.diskRootHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/disk/partitions", Handler: s.diskPartitionHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/disk/disks", Handler: s.diskPhysicalDisksHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/disk/usage", Handler: s.diskUsageHandler, Protected: secure},
		// GPU Routes
		// Network Routes
		{Group: base, HttpMethod: http.MethodGet, Path: "/network", Handler: s.nicRootHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/network/interfaces", Handler: s.nicRootHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/network/interface/:name", Handler: s.nicInterfaceHandler, Protected: secure},
		// Misc Hardware Routes

	}

	// Register routes
	register.RegisterRoutes(routes, s.Controller.SecureMiddleware)
	s.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

// Enabled returns true if this module is enabled otherwise it returns false
func (s *HardwareSubsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the subsystem
func (s *HardwareSubsystem) Name() string {
	return s.name
}
