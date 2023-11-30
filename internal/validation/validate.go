package validation

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/middleware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"go.uber.org/zap"
)

// NewSubsystem creates a new authn subsystem
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "validation"
	vs := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    c.Config.Subsystems.Validation,
		name:       name,
		endpoints:  []*endpoint.Endpoint{},
	}

	return &vs
}

// Subsystem is the subsystem for authentication
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	endpoints  []*endpoint.Endpoint
}

// Handler returns the handler for the middleware
func (s Subsystem) Handler(ep endpoint.Endpoint) endpoint.Func {
	// Bypass authentication for endpoints that don't use auth
	if !s.enabled {
		return ep.Function
	}
	return s.Validate(ep, ep.Function)
}

// Priority returns the priority of the middleware
func (s Subsystem) Priority() middleware.Priority {
	return middleware.PriorityValidation
}

// Validate validates the request and response
func (s Subsystem) Validate(ep endpoint.Endpoint, next endpoint.Func) endpoint.Func {
	return func(in *endpoint.Request) (out *endpoint.Response, err error) {
		s.Logger.Debug("request validation middleware called")
		// Do input validation
		err = ep.ValidateRequest(in)
		if err != nil {
			return nil, err
		}
		return next(in)
		// Do output validation
		//err = ep.ValidateResponse(out)
	}
}

// Endpoints returns the endpoints that this subsystem handles
func (s Subsystem) Endpoints() []*endpoint.Endpoint {
	return s.endpoints
}

// Enabled returns true if the subsystem is enabled
func (s Subsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the subsystem
func (s Subsystem) Name() string {
	return s.name
}
