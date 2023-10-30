package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
)

type Middleware interface {
	Name() string
	Handler(next endpoint.EndpointFunc) endpoint.EndpointFunc
	Priority() MiddlewarePriority
}
