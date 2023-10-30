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
	endpoints  []*endpoint.Endpoint
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
	s.endpoints = []*endpoint.Endpoint{
		// OS Specific Endpoints
		{Path: fmt.Sprintf("%s/", base), Action: endpoint.ActionRead, Function: s.networkInfoHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/arp", base), Action: endpoint.ActionRead, Function: s.arpTableHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/routes", base), Action: endpoint.ActionRead, Function: s.getRoutesHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/route", base), Action: endpoint.ActionRead, Function: s.getRouteHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/route", base), Action: endpoint.ActionWrite, Function: s.updateRouteHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: RouteTableRowArgs{}},
		{Path: fmt.Sprintf("%s/route", base), Action: endpoint.ActionCreate, Function: s.createRouteHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: RouteTableRowArgs{}},
		{Path: fmt.Sprintf("%s/route", base), Action: endpoint.ActionDelete, Function: s.deleteRouteHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: RouteTableRowArgs{}},
		// Unified Endpoints
		{Path: fmt.Sprintf("%s/routes", unified), Action: endpoint.ActionRead, Function: s.getRoutesUnifiedHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/route", unified), Action: endpoint.ActionRead, Function: s.getRouteUnifiedHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/route", unified), Action: endpoint.ActionWrite, Function: s.updateRouteUnifiedHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/route", unified), Action: endpoint.ActionCreate, Function: s.createRouteUnifiedHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
		{Path: fmt.Sprintf("%s/route", unified), Action: endpoint.ActionDelete, Function: s.deleteRouteUnifiedHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil},
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
func (s *Subsystem) Endpoints() []*endpoint.Endpoint {
	return s.endpoints
}

func (s *Subsystem) networkInfoHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return s.NicInfo.Info(), nil
	}, "basic information about the network interfaces on the system")
}
