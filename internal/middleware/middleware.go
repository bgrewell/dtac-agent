package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
)

// Middleware is the interface for middleware
type Middleware interface {
	Name() string
	Handler(ep endpoint.Endpoint) endpoint.Func
	Priority() Priority
}
