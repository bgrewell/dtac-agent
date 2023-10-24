package system

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// NewSubsystem creates a new instance of the Subsystem struct
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "system"
	s := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
	}
	s.info = &Info{}
	s.info.Initialize(s.Logger)
	return &s
}

// Subsystem is a simple example subsystem for showing how the pieces fit together
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	enabled    bool        // Optional subsystems have a boolean to control if they are enabled
	name       string      // Subsystem name
	info       *Info       // Info structure
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
		{Group: base, HTTPMethod: http.MethodGet, Path: "/uuid", Handler: s.uuidHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/product", Handler: s.productHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/os", Handler: s.osHandler, Protected: secure},
	}

	// Register routes
	helpers.RegisterRoutes(routes, s.Controller.SecureMiddleware)
	s.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

// Enabled returns true if the subsystem is enabled
func (s *Subsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the subsystem
func (s *Subsystem) Name() string {
	return s.name
}

func (s *Subsystem) rootHandler(c *gin.Context) {
	start := time.Now()
	s.Controller.Formatter.WriteResponse(c, time.Since(start), s.info)
}

func (s *Subsystem) uuidHandler(c *gin.Context) {
	start := time.Now()
	uuid := s.info.UUID
	s.Controller.Formatter.WriteResponse(c, time.Since(start), uuid)
}

func (s *Subsystem) productHandler(c *gin.Context) {
	start := time.Now()
	product := s.info.ProductName
	s.Controller.Formatter.WriteResponse(c, time.Since(start), product)
}

func (s *Subsystem) osHandler(c *gin.Context) {
	start := time.Now()
	os := s.info.serializeOs()
	s.Controller.Formatter.WriteResponse(c, time.Since(start), os)
}
