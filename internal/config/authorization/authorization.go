package authorization

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"go.uber.org/zap"
	"os"
)

// EnsureAuthzModel ensures the authz subsystem has a default model
func EnsureAuthzModel(c *controller.Controller) {
	// Check if file exists
	if _, err := os.Stat(c.Config.Auth.Model); os.IsNotExist(err) {
		// File doesn't exist, create it with default contents
		defaultContents := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act || g2(r.sub, "admin")
`
		err := os.WriteFile(c.Config.Auth.Model, []byte(defaultContents), 0600)
		if err != nil {
			c.Logger.Error("failed to write authn model file", zap.Error(err))
			return
		}
		c.Logger.Info("created authn model file", zap.String("file", c.Config.Auth.Model))
	}
}

// EnsureAuthzPolicy ensures the authz subsystem has a default policy
func EnsureAuthzPolicy(c *controller.Controller) {
	// Check if file exists
	if _, err := os.Stat(c.Config.Auth.Policy); os.IsNotExist(err) {
		// File doesn't exist, create it with default contents
		defaultContents := `
`
		err := os.WriteFile(c.Config.Auth.Policy, []byte(defaultContents), 0600)
		if err != nil {
			c.Logger.Error("failed to write authn policy file", zap.Error(err))
			return
		}
		c.Logger.Info("created authn policy file", zap.String("file", c.Config.Auth.Policy))
	}
}
