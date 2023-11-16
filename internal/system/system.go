package system

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"go.uber.org/zap"
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
	s.register()
	return &s
}

// Subsystem is a simple example subsystem for showing how the pieces fit together
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	enabled    bool        // Optional subsystems have a boolean to control if they are enabled
	name       string      // Subsystem name
	info       *Info       // Info structure
	endpoints  []*endpoint.Endpoint
}

// register registers the routes that this module handles
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	// Create a group for this subsystem
	base := s.name

	// Endpoints
	secure := s.Controller.Config.Auth.DefaultSecure
	s.endpoints = []*endpoint.Endpoint{
		{Path: fmt.Sprintf("%s/", base), Action: endpoint.ActionRead, Function: s.rootHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil, ExpectedOutput: &Info{}},
		{Path: fmt.Sprintf("%s/uuid", base), Action: endpoint.ActionRead, Function: s.uuidHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil, ExpectedOutput: Info{}.UUID},
		{Path: fmt.Sprintf("%s/product", base), Action: endpoint.ActionRead, Function: s.productHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil, ExpectedOutput: Info{}.ProductName},
		{Path: fmt.Sprintf("%s/os", base), Action: endpoint.ActionRead, Function: s.osHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil, ExpectedOutput: nil},
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

func (s *Subsystem) rootHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return s.info, nil
	}, "system information")
}

func (s *Subsystem) uuidHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return s.info.UUID, nil
	}, "system uuid identifier")
}

func (s *Subsystem) productHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return s.info.ProductName, nil
	}, "system product name")
}

func (s *Subsystem) osHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		return s.info.serializeOs(), nil
	}, "system operation system information")
}
