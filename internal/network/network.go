package network

import (
	"encoding/json"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/shirou/gopsutil/net"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/hardware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
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
	authzAdmin := endpoint.AuthGroupAdmin.String()
	authzUser := endpoint.AuthGroupUser.String()
	s.endpoints = []*endpoint.Endpoint{
		endpoint.NewEndpoint(fmt.Sprintf("%s/", base), endpoint.ActionRead, "network information", s.networkInfoHandler, secure, authzUser, endpoint.WithOutput([]net.InterfaceStat{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/arp", base), endpoint.ActionRead, "arp table information", s.arpTableHandler, secure, authzUser, endpoint.WithOutput([]ArpEntry{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/routes", base), endpoint.ActionRead, "route table information", s.getRoutesHandler, secure, authzUser, endpoint.WithOutput([]RouteTableRow{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", base), endpoint.ActionRead, "route information", s.getRouteHandler, secure, authzUser),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", base), endpoint.ActionWrite, "update existing route", s.updateRouteHandler, secure, authzAdmin, endpoint.WithBody(RouteTableRowArgs{}), endpoint.WithOutput([]RouteTableRow{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", base), endpoint.ActionCreate, "create new route", s.createRouteHandler, secure, authzAdmin, endpoint.WithBody(RouteTableRowArgs{}), endpoint.WithOutput([]RouteTableRow{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", base), endpoint.ActionDelete, "delete route", s.deleteRouteHandler, secure, authzAdmin, endpoint.WithBody(RouteTableRowArgs{}), endpoint.WithOutput([]RouteTableRow{})),

		// Unified Endpoints
		endpoint.NewEndpoint(fmt.Sprintf("%s/routes", unified), endpoint.ActionRead, "os agnostic route table information", s.getRoutesUnifiedHandler, secure, authzUser),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", unified), endpoint.ActionRead, "os agnostic route information", s.getRouteUnifiedHandler, secure, authzUser),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", unified), endpoint.ActionWrite, "os agnostic update exiting route", s.updateRouteUnifiedHandler, secure, authzAdmin),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", unified), endpoint.ActionCreate, "os agnostic create new route", s.createRouteUnifiedHandler, secure, authzAdmin),
		endpoint.NewEndpoint(fmt.Sprintf("%s/route", unified), endpoint.ActionDelete, "os agnostic delete route", s.deleteRouteUnifiedHandler, secure, authzAdmin),
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

func (s *Subsystem) networkInfoHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return json.Marshal(s.NicInfo.Info())
	}, "basic information about the network interfaces on the system")
}
