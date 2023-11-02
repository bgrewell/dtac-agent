package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
)

// AuthorizationMiddleware is the interface for the authorization middleware
type AuthorizationMiddleware interface {
	AuthorizationHandler(next endpoint.Func) endpoint.Func
}
