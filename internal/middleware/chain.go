package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
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
func Chain(middlewares []Middleware, endpoint endpoint.Func) endpoint.Func {
	for _, middleware := range middlewares {
		endpoint = middleware.Handler(endpoint)
	}

	return endpoint
}
