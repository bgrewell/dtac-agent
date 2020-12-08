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

	// PUT Routes - Update information
	r.PUT("/network/route", handlers.UpdateRouteHandler)

	// Delete Routes - Remove information
	r.DELETE("/network/route", handlers.DeleteRouteHandler)

	// POST Routes - Create information
	r.POST("/login", handlers.LoginHandler)
	r.POST("/network/route", handlers.CreateRouteHandler)
}
