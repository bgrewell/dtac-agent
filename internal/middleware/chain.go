package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"sort"
)

// Sort sorts the middleware by priority
func Sort(middlewares []Middleware) []Middleware {
	sort.Slice(middlewares, func(i, j int) bool {
		return middlewares[i].Priority() > middlewares[j].Priority()
	})
	return middlewares
}

// Chain chains the middleware
func Chain(middlewares []Middleware, endpoint endpoint.Endpoint) endpoint.Func {
	for _, middleware := range middlewares {
		endpoint.Function = middleware.Handler(endpoint)
	}

	return endpoint.Function
}
