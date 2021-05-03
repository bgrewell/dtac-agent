package httprouting

import (
	"github.com/BGrewell/system-api/handlers"
	"github.com/gin-gonic/gin"
)

func AddGeneralHandlers(r *gin.Engine) {
	// GET Routes - Retrieve information
	r.GET("/", handlers.HomeHandler)
	r.GET("/network/interfaces", handlers.GetInterfacesHandler)
	r.GET("/network/interfaces/names", handlers.GetInterfaceNamesHandler)
	r.GET("/network/interface/:name", handlers.GetInterfaceByNameHandler)
	r.GET("/network/routes", handlers.GetRoutesHandler)
	r.GET("/secret", handlers.SecretTestHandler)
	r.GET("/iperf/client/results/:id", handlers.GetIperfClientTestResultsHandler)
	r.GET("/iperf/client/live/:id", handlers.GetIperfClientTestLiveHandler)
	r.GET("/iperf/server/results/:id", handlers.GetIperfServerTestResultsHandler)

	// PUT Routes - Update information
	r.PUT("/network/route", handlers.UpdateRouteHandler)

	// Delete Routes - Remove information
	r.DELETE("/network/route", handlers.DeleteRouteHandler)
	r.DELETE("/iperf/client/stop/:id", handlers.DeleteIperfClientTestHandler)
	r.DELETE("/iperf/server/stop/:id", handlers.DeleteIperfServerTestHandler)
	r.DELETE("/iperf/reset", handlers.DeleteIperfResetHandler)

	// POST Routes - Create information
	r.POST("/login", handlers.LoginHandler)
	r.POST("/network/route", handlers.CreateRouteHandler)
	r.POST("/iperf/client/start/:host", handlers.CreateIperfClientTestHandler)
	r.POST("/iperf/server/start", handlers.CreateIperfServerTestHandler)
	r.POST("/iperf/server/start/:bind", handlers.CreateIperfServerTestHandler)

}
