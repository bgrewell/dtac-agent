package system

import (
	"encoding/json"
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
	authz := endpoint.AuthGroupUser.String()
	s.endpoints = []*endpoint.Endpoint{
		endpoint.NewEndpoint(fmt.Sprintf("%s/", base), endpoint.ActionRead, "general system information", s.rootHandler, secure, authz, endpoint.WithOutput(&Info{})),
		endpoint.NewEndpoint(fmt.Sprintf("%s/uuid", base), endpoint.ActionRead, "system uuid", s.uuidHandler, secure, authz, endpoint.WithOutput(Info{}.UUID)),
		endpoint.NewEndpoint(fmt.Sprintf("%s/product", base), endpoint.ActionRead, "system product", s.productHandler, secure, authz, endpoint.WithOutput(Info{}.ProductName)),
		endpoint.NewEndpoint(fmt.Sprintf("%s/os", base), endpoint.ActionRead, "operating system", s.osHandler, secure, authz),
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

func (s *Subsystem) rootHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return json.Marshal(s.info)
	}, "system information")
}

func (s *Subsystem) uuidHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return json.Marshal(s.info.UUID)
	}, "system uuid identifier")
}

func (s *Subsystem) productHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return json.Marshal(s.info.ProductName)
	}, "system product name")
}

func (s *Subsystem) osHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		return json.Marshal(s.info.serializeOs())
	}, "system operation system information")
}
