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

// NewAuthzSubsystem creates a new authz subsystem
func NewAuthzSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "authz"

	az := AuthzSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    true,
		name:       name,
		enforcer:   nil,
	}
	return &az
}

// AuthzSubsystem is the subsystem for authorization
type AuthzSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string
	enforcer   *casbin.Enforcer
}

// Register registers the authz subsystem
func (as *AuthzSubsystem) Register() error {
	if !as.Enabled() {
		as.Logger.Info("subsystem is disabled", zap.String("subsystem", as.Name()))
		return nil
	}

	enforcer, err := casbin.NewEnforcer(config.DEFAULT_AUTH_MODEL_NAME, config.DEFAULT_AUTH_POLICY_NAME)
	if err != nil {
		as.Logger.Fatal("failed to create casbin enforcer", zap.Error(err))
	}
	as.enforcer = enforcer

	return nil
}

// Enabled returns whether or not the authz subsystem is enabled
func (as *AuthzSubsystem) Enabled() bool {
	return as.enabled
}

// AuthorizationHandler is the handler for authorization
func (as *AuthzSubsystem) AuthorizationHandler(c *gin.Context) {
	// This is just a extremely basic authorization function right now. Will need to be built out to have full
	// RBAC or ACL access controls in place. This implementation just checks to see if the user can access the
	// resource and the default model says that "admin" can access anything.
	if user, ok := c.Get("username"); ok {
		if res, _ := as.enforcer.Enforce(user, c.Request.URL.Path, c.Request.Method); res {
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
func (as *AuthzSubsystem) Name() string {
	return as.name
}
