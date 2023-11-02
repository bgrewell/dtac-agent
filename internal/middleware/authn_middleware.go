package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
)

// AuthenticationMiddleware is the interface for the authentication middleware
type AuthenticationMiddleware interface {
	AuthenticationHandler(next endpoint.Func) endpoint.Func
}
