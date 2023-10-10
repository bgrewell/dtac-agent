package system

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/register"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "system"
	s := SystemSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
	}
	s.info = &SystemInfo{}
	s.info.Initialize(s.Logger)
	return &s
}

// SystemSubsystem is a simple example subsystem for showing how the pieces fit together
type SystemSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	enabled    bool        // Optional subsystems have a boolean to control if they are enabled
	name       string      // Subsystem name
	info       *SystemInfo // SystemInfo structure
}

// Register() registers the routes that this module handles
func (s *SystemSubsystem) Register() error {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return nil
	}

	// Create a group for this subsystem
	base := s.Controller.Router.Group(s.name)

	// Routes
	secure := s.Controller.Config.Auth.DefaultSecure
	routes := []types.RouteInfo{
		{Group: base, HttpMethod: http.MethodGet, Path: "/", Handler: s.rootHandler, Protected: secure},
	}

	// Register routes
	register.RegisterRoutes(routes, s.Controller.SecureMiddleware)
	s.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

func (s *SystemSubsystem) Enabled() bool {
	return s.enabled
}

func (s *SystemSubsystem) Name() string {
	return s.name
}

func (s *SystemSubsystem) rootHandler(c *gin.Context) {
	start := time.Now()
	helpers.WriteResponseJSON(c, time.Since(start), s.info)
}
