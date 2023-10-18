package http

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
)

// NewRouteList creates a new instance of the RouteList struct
func NewRouteList(router *gin.Engine, cfg *config.Configuration, log *zap.Logger) *RouteList {
	httpList := RouteList{
		Router: router,
		Config: cfg,
		Logger: log.With(zap.String("module", "route_list")),
	}
	httpList.UpdateRoutes()
	return &httpList
}

// RouteList is the struct for the http route list
type RouteList struct {
	Routes []*RouteInfo `json:"routes"`
	Router *gin.Engine
	Config *config.Configuration
	Logger *zap.Logger
}

// UpdateRoutes updates the http route list
func (hrl *RouteList) UpdateRoutes() {
	routes := hrl.Router.Routes()
	hrl.Routes = make([]*RouteInfo, len(routes))
	for idx, route := range routes {
		hrl.Routes[idx] = &RouteInfo{
			Method:  route.Method,
			Path:    route.Path,
			Handler: route.Handler,
		}
	}
}

// RouteInfo is the struct for the http route info
type RouteInfo struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"-"`
}
