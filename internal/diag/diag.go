package diag

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version"
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
	ds.register()
	return &ds
}

// Subsystem is the subsystem that contains routes related to internal dtac diagnostics
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string // Subsystem name
	endpoints  []endpoint.Endpoint
}

// register registers the endpoints that this module handles
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	// Create a group for this subsystem
	base := s.name

	// Endpoints
	secure := s.Controller.Config.Auth.DefaultSecure
	s.endpoints = []endpoint.Endpoint{
		{fmt.Sprintf("%s/", base), endpoint.ActionRead, s.rootHandler, secure, nil, nil},
		{fmt.Sprintf("%s/jwt", base), endpoint.ActionRead, s.jwtTestHandler, secure, nil, nil},
		{fmt.Sprintf("%s/endpoints", base), endpoint.ActionRead, s.endpointListPrintHandler, secure, nil, nil},
		{fmt.Sprintf("%s/runningas", base), endpoint.ActionRead, s.runningAsHandler, secure, nil, nil},
	}
}

// Enabled returns true if this module is enabled otherwise it returns false
func (s *Subsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the Subsystem
func (s *Subsystem) Name() string {
	return s.name
}

// Endpoints returns an array of endpoints that this Subsystem handles
func (s *Subsystem) Endpoints() []endpoint.Endpoint {
	return s.endpoints
}

// rootHandler handles requests for the root path for this subsystem
func (s *Subsystem) rootHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return version.Current(), nil
	}, "diagnostic information")
}

// endpointListPrintHandler handles requests for the supported endpoints
func (s *Subsystem) endpointListPrintHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return s.Controller.EndpointList, nil
	}, "enabled api endpoints")
}

// jwtTestHandler is a legacy endpoint that is depreciated and needs to be removed
func (s *Subsystem) jwtTestHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return map[string]string{"message": "jwt test page"}, nil
	}, "jwt authentication test")
}

// runningAsHandler returns information about the user and group context the application is running as
func (s *Subsystem) runningAsHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return AgentRunningAsUser()
	}, "application running as user/group information")
}
