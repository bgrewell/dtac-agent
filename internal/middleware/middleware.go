package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
)

// Middleware is the interface for middleware
type Middleware interface {
	Name() string
	Handler(next endpoint.Func) endpoint.Func
	Priority() Priority
}
