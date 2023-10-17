package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/hardware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/register"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "network"
	ns := NetworkSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    c.Config.Subsystems.Network,
		name:       name,
	}
	return &ns
}

// NetworkSubsystem handles network related functionalities
type NetworkSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	NicInfo    hardware.NicInfo
	enabled    bool   // Optional subsystems have a boolean to control if they are enabled
	name       string // Subsystem name
}

// Register() registers the routes that this module handles. Currently empty as no routes defined.
func (s *NetworkSubsystem) Register() error {
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
		{Group: base, HttpMethod: http.MethodGet, Path: "/", Handler: s.networkInfoHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/arp", Handler: arpTableHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/routes", Handler: getRoutesHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/route", Handler: getRouteHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodPut, Path: "/route", Handler: updateRouteHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodPost, Path: "/route", Handler: createRouteHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodDelete, Path: "/route", Handler: deleteRouteHandler, Protected: secure},
		{Group: unified, HttpMethod: http.MethodGet, Path: "/routes", Handler: getRoutesUnifiedHandler, Protected: secure},
		{Group: unified, HttpMethod: http.MethodGet, Path: "/route", Handler: getRouteUnifiedHandler, Protected: secure},
		{Group: unified, HttpMethod: http.MethodPut, Path: "/route", Handler: updateRouteUnifiedHandler, Protected: secure},
		{Group: unified, HttpMethod: http.MethodPost, Path: "/route", Handler: createRouteUnifiedHandler, Protected: secure},
		{Group: unified, HttpMethod: http.MethodDelete, Path: "/route", Handler: deleteRouteUnifiedHandler, Protected: secure},
	}

	// Register routes
	register.RegisterRoutes(routes, s.Controller.SecureMiddleware)
	s.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

func (s *NetworkSubsystem) Enabled() bool {
	return s.enabled
}

func (s *NetworkSubsystem) Name() string {
	return s.name
}

func (s *NetworkSubsystem) networkInfoHandler(c *gin.Context) {
	start := time.Now()
	s.NicInfo.Update()
	response := gin.H{
		"network-interfaces": types.AnnotatedStruct{
			Description: "returns basic information about the network interfaces on the system",
			Value:       s.NicInfo.Info(),
		},
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}
