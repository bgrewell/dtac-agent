package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
)

// RegisterRoutes registers the routes that this module handles
func RegisterRoutes(routes []types.RouteInfo, secureMw []gin.HandlerFunc) {
	for _, route := range routes {
		funcs := []gin.HandlerFunc{}
		if route.Protected {
			funcs = append(funcs, secureMw...)
		}
		funcs = append(funcs, route.Handler)
		route.Group.Handle(route.HTTPMethod, route.Path, funcs...)
	}
}
