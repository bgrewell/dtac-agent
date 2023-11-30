package authz

import (
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

// CasbinLogger is an implementation of casbin.Logger interface
type CasbinLogger struct {
	enabled bool
	logger  *zap.Logger
}

// EnableLog controls whether to enable the logger
func (l *CasbinLogger) EnableLog(enable bool) {
	l.enabled = enable
}

// IsEnabled returns true if the logger is enabled
func (l *CasbinLogger) IsEnabled() bool {
	return l.enabled
}

// LogModel logs model information, only called when the model changes
func (l *CasbinLogger) LogModel(model [][]string) {
	if !l.enabled {
		return
	}
	l.logger.Info("logging model", zap.Any("model", model))
}

// LogEnforce logs the enforcing information, including the request, the policy, and the decision
func (l *CasbinLogger) LogEnforce(matcher string, request []interface{}, result bool, explains [][]string) {
	if !l.enabled {
		return
	}
	l.logger.Info("logging enforcement", zap.String("matcher", matcher), zap.Any("request", request), zap.Bool("result", result), zap.Any("explains", explains))
}

// LogRole logs the role information, including role, users, and permissions
func (l *CasbinLogger) LogRole(roles []string) {
	if !l.enabled {
		return
	}
	l.logger.Info("logging roles", zap.Strings("roles", roles))
}

// LogPolicy logs the policy information, including policy type, role, and permission
func (l *CasbinLogger) LogPolicy(policy map[string][][]string) {
	if !l.enabled {
		return
	}

	l.logger.Info("logging policy", zap.Any("policy", policy))
}

// LogError logs error message
func (l *CasbinLogger) LogError(err error, msg ...string) {
	if !l.enabled {
		return
	}

	if len(msg) > 0 {
		l.logger.Error("casbin error", zap.Error(err), zap.Strings("msg", msg))
	} else {
		l.logger.Error("casbin error", zap.Error(err))
	}
}

// LogCurrentPolicies logs current policies
func (l *CasbinLogger) LogCurrentPolicies(enforcer *casbin.Enforcer) {
	if !l.enabled {
		return
	}

	l.logger.Info("Current Policies:")
	for _, policy := range enforcer.GetPolicy() {
		l.logger.Info("Policy:", zap.Any("policy", policy))
	}

	l.logger.Info("Current Grouping Policies:")
	for _, gPolicy := range enforcer.GetGroupingPolicy() {
		l.logger.Info("Grouping Policy:", zap.Any("grouping policy", gPolicy))
	}

	l.logger.Info("Current Role Hierarchies:")
	for _, hierarchy := range enforcer.GetNamedGroupingPolicy("g2") {
		l.logger.Info("Role Hierarchy:", zap.Any("role hierarchy", hierarchy))
	}
}
