package middleware

import (
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// Middleware is the interface for middleware
type Middleware interface {
	Name() string
	Handler(ep endpoint.Endpoint) endpoint.Func
	Priority() Priority
}
