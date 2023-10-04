package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
)

func NewHttpRouteList(router *gin.Engine, cfg *config.Configuration, log *zap.Logger) *HttpRouteList {
	httpList := HttpRouteList{
		Router: router,
		Config: cfg,
		Logger: log.With(zap.String("module", "route_list")),
	}
	httpList.UpdateRoutes()
	return &httpList
}

type HttpRouteList struct {
	Routes []*HttpRouteInfo `json:"routes"`
	Router *gin.Engine
	Config *config.Configuration
	Logger *zap.Logger
}

func (hrl *HttpRouteList) UpdateRoutes() {
	routes := hrl.Router.Routes()
	hrl.Routes = make([]*HttpRouteInfo, len(routes))
	for idx, route := range routes {
		hrl.Routes[idx] = &HttpRouteInfo{
			Method:  route.Method,
			Path:    route.Path,
			Handler: route.Handler,
		}
	}
}

func (hrl *HttpRouteList) HttpRoutePrintHandler(c *gin.Context) {
	hrl.UpdateRoutes()
	c.IndentedJSON(http.StatusOK, gin.H{
		"routes": types.AnnotatedStruct{
			Description: "list of registered http endpoints being served by dtac",
			Value:       hrl.Routes,
		},
	})
}

type HttpRouteInfo struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"-"`
}
