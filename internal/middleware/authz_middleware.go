package middleware

import (
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// AuthorizationMiddleware is the interface for the authorization middleware
type AuthorizationMiddleware interface {
	AuthorizationHandler(next endpoint.Func) endpoint.Func
	RegisterPolicies() error
}
