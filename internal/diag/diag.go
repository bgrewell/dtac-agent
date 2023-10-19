package diag

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
)

// NewSubsystem creates a new instances of the Subsystem and if that subsystem is enabled it calls
// the Register() function to register the routes that the Subsystem handles
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "diag"
	ds := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    c.Config.Subsystems.Diag,
		name:       name,
	}
	return &ds
}

// Subsystem is the subsystem that contains routes related to internal dtac diagnostics
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string // Subsystem name
}

// Register registers the routes that this module handles
func (s *Subsystem) Register() error {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return nil
	}
	// Create a group for this subsystem
	base := s.Controller.Router.Group(s.name)

	// Routes
	secure := s.Controller.Config.Auth.DefaultSecure
	routes := []types.RouteInfo{
		{Group: base, HTTPMethod: http.MethodGet, Path: "/", Handler: s.rootHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/jwt", Handler: s.jwtTestHandler, Protected: true},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/routes", Handler: s.httpRoutePrintHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/runningas", Handler: s.runningAsHandler, Protected: false},
	}

	// Register routes
	helpers.RegisterRoutes(routes, s.Controller.SecureMiddleware)
	s.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

// Enabled returns true if this module is enabled otherwise it returns false
func (s *Subsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the Subsystem
func (s *Subsystem) Name() string {
	return s.name
}

// rootHandler handles requests for the root path for this subsystem
func (s *Subsystem) rootHandler(c *gin.Context) {
	start := time.Now()
	response := gin.H{
		"version": types.AnnotatedStruct{
			Description: fmt.Sprintf("%s version information", s.Controller.Config.Internal.ShortName),
			Value:       version.Current(),
		},
		"memory": types.AnnotatedStruct{
			Description: fmt.Sprintf("current %s memory usage", s.Controller.Config.Internal.ShortName),
			Value:       CurrentMemoryStats(),
		},
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}

func (s *Subsystem) httpRoutePrintHandler(c *gin.Context) {
	start := time.Now()
	s.Controller.HTTPRouteList.UpdateRoutes()
	response := gin.H{
		"routes": types.AnnotatedStruct{
			Description: fmt.Sprintf("list of registered http endpoints being served by %s", s.Controller.Config.Internal.ShortName),
			Value:       s.Controller.HTTPRouteList.Routes,
		},
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}

func (s *Subsystem) jwtTestHandler(c *gin.Context) {
	start := time.Now()
	response := gin.H{
		"message": "jwt test page",
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}

func (s *Subsystem) runningAsHandler(c *gin.Context) {
	start := time.Now()
	ug, err := AgentRunningAsUser()
	if err != nil {
		helpers.WriteErrorResponseJSON(c, err)
		return
	}
	response := gin.H{
		"runningAs": ug,
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}
