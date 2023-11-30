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

func (s Subsystem) Handler(ep endpoint.Endpoint) endpoint.Func {
	// Bypass authentication for endpoints that don't use auth
	if !s.enabled {
		return ep.Function
	}
	return s.Validate(ep, ep.Function)
}

func (s Subsystem) Priority() middleware.Priority {
	return middleware.PriorityValidation
}

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

func (s Subsystem) Endpoints() []*endpoint.Endpoint {
	return s.endpoints
}

func (s Subsystem) Enabled() bool {
	return s.enabled
}

func (s Subsystem) Name() string {
	return s.name
}
