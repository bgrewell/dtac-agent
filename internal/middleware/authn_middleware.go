package middleware

import (
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// AuthenticationMiddleware is the interface for the authentication middleware
type AuthenticationMiddleware interface {
	AuthenticationHandler(next endpoint.Func) endpoint.Func
}
