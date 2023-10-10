package hardware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/register"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
)

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
	}
	return &hw
}

type HardwareSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string // Subsystem name
	nic        NicInfo
	cpu        CpuInfo
	mem        MemoryInfo
}

// Register() registers the routes that this module handles
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

func (s *HardwareSubsystem) Name() string {
	return s.name
}
