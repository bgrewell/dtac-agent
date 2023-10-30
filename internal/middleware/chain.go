package middleware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"sort"
)

func Sort(middlewares []Middleware) []Middleware {
	sort.Slice(middlewares, func(i, j int) bool {
		return middlewares[i].Priority() > middlewares[j].Priority()
	})
	return middlewares
}

func Chain(middlewares []Middleware, endpoint endpoint.EndpointFunc) endpoint.EndpointFunc {
	for _, middleware := range middlewares {
		endpoint = middleware.Handler(endpoint)
	}

	return endpoint
}
