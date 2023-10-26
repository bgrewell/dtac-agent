package network

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/hardware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"go.uber.org/zap"
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
	ns.register()
	return &ns
}

// Subsystem handles network related functionalities
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	NicInfo    hardware.NicInfo
	enabled    bool   // Optional subsystems have a boolean to control if they are enabled
	name       string // Subsystem name
	endpoints  []endpoint.Endpoint
}

// register registers the routes that this module handles. Currently empty as no routes defined.
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	// Create group(s) for this subsystem
	base := s.name
	unified := fmt.Sprintf("u/%s", s.name)

	// Endpoints
	secure := s.Controller.Config.Auth.DefaultSecure
	s.endpoints = []endpoint.Endpoint{
		// OS Specific Endpoints
		{fmt.Sprintf("%s/", base), endpoint.ActionRead, s.networkInfoHandler, secure, nil, nil},
		{fmt.Sprintf("%s/arp", base), endpoint.ActionRead, s.arpTableHandler, secure, nil, nil},
		{fmt.Sprintf("%s/routes", base), endpoint.ActionRead, s.getRoutesHandler, secure, nil, nil},
		{fmt.Sprintf("%s/route", base), endpoint.ActionRead, s.getRouteHandler, secure, nil, nil},
		{fmt.Sprintf("%s/route", base), endpoint.ActionWrite, s.updateRouteHandler, secure, nil, RouteTableRow{}},
		{fmt.Sprintf("%s/route", base), endpoint.ActionCreate, s.createRouteHandler, secure, nil, RouteTableRow{}},
		{fmt.Sprintf("%s/route", base), endpoint.ActionDelete, s.deleteRouteHandler, secure, nil, RouteTableRow{}},
		// Unified Endpoints
		{fmt.Sprintf("%s/routes", unified), endpoint.ActionRead, s.getRoutesUnifiedHandler, secure, nil, nil},
		{fmt.Sprintf("%s/route", unified), endpoint.ActionRead, s.getRouteUnifiedHandler, secure, nil, nil},
		{fmt.Sprintf("%s/route", unified), endpoint.ActionWrite, s.updateRouteUnifiedHandler, secure, nil, nil},
		{fmt.Sprintf("%s/route", unified), endpoint.ActionCreate, s.createRouteUnifiedHandler, secure, nil, nil},
		{fmt.Sprintf("%s/route", unified), endpoint.ActionDelete, s.deleteRouteUnifiedHandler, secure, nil, nil},
	}
}

// Enabled returns true if the subsystem is enabled
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

func (s *Subsystem) networkInfoHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return s.NicInfo.Info(), nil
	}, "basic information about the network interfaces on the system")
}
