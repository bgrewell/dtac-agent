package register

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
)

func RegisterRoutes(routes []types.RouteInfo, secureMw []gin.HandlerFunc) {
	for _, route := range routes {
		funcs := []gin.HandlerFunc{}
		if route.Protected {
			funcs = append(funcs, secureMw...)
		}
		funcs = append(funcs, route.Handler)
		route.Group.Handle(route.HttpMethod, route.Path, funcs...)
	}
}
