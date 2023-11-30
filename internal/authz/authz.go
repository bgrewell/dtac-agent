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
	"strings"
)

// NewSubsystem creates a new authz subsystem
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "authz"

	az := Subsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    c.Config.Subsystems.Auth,
		name:       name,
		enforcer:   nil,
	}
	az.register()
	return &az
}

// Subsystem is the subsystem for authorization
type Subsystem struct {
	Controller   *controller.Controller
	Logger       *zap.Logger
	enabled      bool
	name         string
	enforcer     *casbin.Enforcer
	policyLogger *CasbinLogger
	endpoints    []*endpoint.Endpoint
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

	s.policyLogger = &CasbinLogger{
		enabled: true,
		logger:  s.Logger.With(zap.String("module", "casbin")),
	}
	s.enforcer.EnableLog(true)
	s.enforcer.SetLogger(s.policyLogger)
	// Setup role hierarchy
	addRoleHierarchies(enforcer, s.Logger)
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
	if !ep.Secure || !s.enabled {
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
	return func(in *endpoint.Request) (out *endpoint.Response, err error) {
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
		userJSON := in.Metadata[types.ContextAuthUser.String()]
		err = json.Unmarshal([]byte(userJSON), &user)
		if err != nil {
			return nil, fmt.Errorf("error retrieving user from context: %v", err)
		}

		// Extract action and path from context
		action := in.Metadata[types.ContextResourceAction.String()]
		path := in.Metadata[types.ContextResourcePath.String()]

		s.Logger.Debug("Username", zap.String("username", user.Username))

		// Get users roles
		roles, err := s.enforcer.GetRolesForUser(user.Username)
		if err != nil {
			return nil, fmt.Errorf("error retrieving roles for user: %v", err)
		}

		// Check if user has access to the resource
		for _, role := range roles {
			canAccess, err := s.enforcer.Enforce(role, path, action)
			if err != nil {
				return nil, fmt.Errorf("error checking role access: %v", err)
			}
			if canAccess {
				in.Metadata[types.ContextAuthRoles.String()] = strings.Join(roles, ",")
				return next(in)
			}
		}

		return nil, errors.New("user not authorized to access this resource")
	}
}

// RegisterPolicies registers the policies for the authz subsystem
func (s *Subsystem) RegisterPolicies() error {
	// Setup policy assignments for users
	users, err := s.Controller.AuthDB.ViewUsers()
	if err != nil {
		return fmt.Errorf("failed to view users: %v", err)
	}
	for _, user := range users {
		for _, group := range user.Groups {
			_, err = s.enforcer.AddGroupingPolicy(user.Username, group)
			if err != nil {
				return err
			}
		}
	}

	// Setup policies for the endpoints
	for _, endpoint := range s.Controller.EndpointList.Endpoints {
		_, err = s.enforcer.AddPolicy(endpoint.AuthGroup, endpoint.Path, endpoint.Action.String())
		if err != nil {
			return err
		}
	}

	s.policyLogger.LogCurrentPolicies(s.enforcer)

	return nil
}

func addRoleHierarchies(enforcer *casbin.Enforcer, logger *zap.Logger) {
	roleHierarchies := []struct {
		parent string
		child  string
	}{
		{"admin", "operator"},
		{"operator", "user"},
		{"user", "guest"},
	}

	for _, hierarchy := range roleHierarchies {
		if _, err := enforcer.AddNamedGroupingPolicy("g2", hierarchy.parent, hierarchy.child); err != nil {
			logger.Fatal("failed to add role hierarchy", zap.String("parent", hierarchy.parent), zap.String("child", hierarchy.child), zap.Error(err))
		}
	}
}
