package middleware

import (
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// ValidationMiddleware is the interface for the authentication middleware
type ValidationMiddleware interface {
	Validate(endpoint *endpoint.Endpoint, next endpoint.Func) endpoint.Func
}
