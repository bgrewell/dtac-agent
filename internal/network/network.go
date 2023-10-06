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
func (ns *NetworkSubsystem) Register() error {
	if !ns.Enabled() {
		ns.Logger.Info("subsystem is disabled", zap.String("subsystem", ns.Name()))
		return nil
	}

	// Create a group for this subsystem
	base := ns.Controller.Router.Group(ns.name)
	unified := ns.Controller.Router.Group(fmt.Sprintf("u/%s", ns.name))

	// Routes
	routes := []types.RouteInfo{
		{Group: base, HttpMethod: http.MethodGet, Path: "/", Handler: ns.networkInfoHandler, Protected: false},
		{Group: base, HttpMethod: http.MethodGet, Path: "/arp", Handler: arpTableHandler, Protected: false},
		{Group: base, HttpMethod: http.MethodGet, Path: "/routes", Handler: getRoutesHandler, Protected: false},
		{Group: base, HttpMethod: http.MethodGet, Path: "/route", Handler: getRouteHandler, Protected: false},
		{Group: base, HttpMethod: http.MethodPut, Path: "/route", Handler: updateRouteHandler, Protected: false},
		{Group: base, HttpMethod: http.MethodPost, Path: "/route", Handler: createRouteHandler, Protected: false},
		{Group: base, HttpMethod: http.MethodDelete, Path: "/route", Handler: deleteRouteHandler, Protected: false},
		{Group: unified, HttpMethod: http.MethodGet, Path: "/routes", Handler: getRoutesUnifiedHandler, Protected: false},
		{Group: unified, HttpMethod: http.MethodGet, Path: "/route", Handler: getRouteUnifiedHandler, Protected: false},
		{Group: unified, HttpMethod: http.MethodPut, Path: "/route", Handler: updateRouteUnifiedHandler, Protected: false},
		{Group: unified, HttpMethod: http.MethodPost, Path: "/route", Handler: createRouteUnifiedHandler, Protected: false},
		{Group: unified, HttpMethod: http.MethodDelete, Path: "/route", Handler: deleteRouteUnifiedHandler, Protected: false},
	}

	// Register routes
	register.RegisterRoutes(routes, ns.Controller.SecureMiddleware)
	ns.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

func (ns *NetworkSubsystem) Enabled() bool {
	return ns.enabled
}

func (ns *NetworkSubsystem) Name() string {
	return ns.name
}

func (ns *NetworkSubsystem) networkInfoHandler(c *gin.Context) {
	start := time.Now()
	ns.NicInfo.Update()
	response := gin.H{
		"network-interfaces": types.AnnotatedStruct{
			Description: "returns basic information about the network interfaces on the system",
			Value:       ns.NicInfo.Info(),
		},
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}
