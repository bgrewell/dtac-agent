package authz

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authndb"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/middleware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"go.uber.org/zap"
)

// NewSubsystem creates a new authz subsystem
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "authz"

	az := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
		enforcer:   nil,
	}
	az.register()
	return &az
}

// Subsystem is the subsystem for authorization
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	enforcer   *casbin.Enforcer
	endpoints  []*endpoint.Endpoint
}

// register registers the authz subsystem
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	enforcer, err := casbin.NewEnforcer(config.DefaultAuthModelName, config.DefaultAuthPolicyName)
	if err != nil {
		s.Logger.Fatal("failed to create casbin enforcer", zap.Error(err))
	}
	s.enforcer = enforcer
}

// Enabled returns whether the authz subsystem is enabled
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

// Handler handles the authentication middleware
func (s *Subsystem) Handler(ep endpoint.Endpoint) endpoint.Func {
	// Bypass authentication for endpoints that don't use auth
	if !ep.Secure {
		return ep.Function
	}
	return s.AuthorizationHandler(ep.Function)
}

// Priority returns the priority of the middleware
func (s *Subsystem) Priority() middleware.Priority {
	return middleware.PriorityAuthorization
}

// AuthorizationHandler is the handler for authorization
func (s *Subsystem) AuthorizationHandler(next endpoint.Func) endpoint.Func {
	return func(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
		s.Logger.Debug("authorization middleware called")

		// Check for metadata that is needed for authorization
		if _, ok := in.Metadata[types.ContextAuthUser.String()]; !ok {
			return nil, errors.New("user is not logged in")
		}
		if _, ok := in.Metadata[types.ContextResourceAction.String()]; !ok {
			return nil, errors.New("resource action is not specified")
		}
		if _, ok := in.Metadata[types.ContextResourcePath.String()]; !ok {
			return nil, errors.New("resource path is not specified")
		}

		var user authndb.User
		userJson := in.Metadata[types.ContextAuthUser.String()]
		err = json.Unmarshal([]byte(userJson), &user)
		if err != nil {
			return nil, fmt.Errorf("error retrieving user from context: %v", err)
		}

		// Extract action and path from context
		action := in.Metadata[types.ContextResourceAction.String()]
		path := in.Metadata[types.ContextResourcePath.String()]

		s.Logger.Debug("Username", zap.String("username", user.Username))

		if canAccess, _ := s.enforcer.Enforce(user.Username, path, string(action)); canAccess {
			return next(in)
		}

		return nil, errors.New("user not authorized to access this resource")
	}
	// This is just a extremely basic authorization function right now. Will need to be built out to have full
	// RBAC or ACL access controls in place. This implementation just checks to see if the user can access the
	// resource and the default model says that "admin" can access anything.
}
