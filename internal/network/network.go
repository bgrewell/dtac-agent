package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/hardware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// NewSubsystem creates a new instance of the Subsystem struct
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "network"
	ns := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    c.Config.Subsystems.Network,
		name:       name,
	}
	return &ns
}

// Subsystem handles network related functionalities
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	NicInfo    hardware.NicInfo
	enabled    bool   // Optional subsystems have a boolean to control if they are enabled
	name       string // Subsystem name
}

// Register registers the routes that this module handles. Currently empty as no routes defined.
func (s *Subsystem) Register() error {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return nil
	}

	// Create a group for this subsystem
	base := s.Controller.Router.Group(s.name)
	unified := s.Controller.Router.Group(fmt.Sprintf("u/%s", s.name))

	// Routes
	secure := s.Controller.Config.Auth.DefaultSecure
	routes := []types.RouteInfo{
		{Group: base, HTTPMethod: http.MethodGet, Path: "/", Handler: s.networkInfoHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/arp", Handler: s.arpTableHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/routes", Handler: s.getRoutesHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodGet, Path: "/route", Handler: s.getRouteHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodPut, Path: "/route", Handler: s.updateRouteHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodPost, Path: "/route", Handler: s.createRouteHandler, Protected: secure},
		{Group: base, HTTPMethod: http.MethodDelete, Path: "/route", Handler: s.deleteRouteHandler, Protected: secure},
		{Group: unified, HTTPMethod: http.MethodGet, Path: "/routes", Handler: s.getRoutesUnifiedHandler, Protected: secure},
		{Group: unified, HTTPMethod: http.MethodGet, Path: "/route", Handler: s.getRouteUnifiedHandler, Protected: secure},
		{Group: unified, HTTPMethod: http.MethodPut, Path: "/route", Handler: s.updateRouteUnifiedHandler, Protected: secure},
		{Group: unified, HTTPMethod: http.MethodPost, Path: "/route", Handler: s.createRouteUnifiedHandler, Protected: secure},
		{Group: unified, HTTPMethod: http.MethodDelete, Path: "/route", Handler: s.deleteRouteUnifiedHandler, Protected: secure},
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

func (s *Subsystem) networkInfoHandler(c *gin.Context) {
	start := time.Now()
	s.NicInfo.Update()
	response := gin.H{
		"network-interfaces": types.AnnotatedStruct{
			Description: "returns basic information about the network interfaces on the system",
			Value:       s.NicInfo.Info(),
		},
	}
	s.Controller.Formatter.WriteResponse(c, time.Since(start), response)
}
