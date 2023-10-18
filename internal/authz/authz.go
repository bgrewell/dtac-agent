package authz

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
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
	return &az
}

// Subsystem is the subsystem for authorization
type Subsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	enforcer   *casbin.Enforcer
}

// Register registers the authz subsystem
func (s *Subsystem) Register() error {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return nil
	}

	enforcer, err := casbin.NewEnforcer(config.DefaultAuthModelName, config.DefaultAuthPolicyName)
	if err != nil {
		s.Logger.Fatal("failed to create casbin enforcer", zap.Error(err))
	}
	s.enforcer = enforcer

	return nil
}

// Enabled returns whether or not the authz subsystem is enabled
func (s *Subsystem) Enabled() bool {
	return s.enabled
}

// AuthorizationHandler is the handler for authorization
func (s *Subsystem) AuthorizationHandler(c *gin.Context) {
	// This is just a extremely basic authorization function right now. Will need to be built out to have full
	// RBAC or ACL access controls in place. This implementation just checks to see if the user can access the
	// resource and the default model says that "admin" can access anything.
	if user, ok := c.Get("username"); ok {
		if res, _ := s.enforcer.Enforce(user, c.Request.URL.Path, c.Request.Method); res {
			c.Next()
		} else {
			helpers.WriteUnauthorizedResponseJSON(c, errors.New("user not authorized to access this resource"))
			return
		}
	} else {
		helpers.WriteUnauthorizedResponseJSON(c, errors.New("user is not logged in"))
		return
	}
}

// Name returns the name of the subsystem
func (s *Subsystem) Name() string {
	return s.name
}
