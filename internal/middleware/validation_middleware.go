package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
)

// ValidationMiddleware is the interface for the authentication middleware
type ValidationMiddleware interface {
	Validate(endpoint *endpoint.Endpoint, next endpoint.Func) endpoint.Func
}
