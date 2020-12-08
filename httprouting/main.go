package httprouting

import (
	"github.com/BGrewell/system-api/handlers"
	"github.com/gin-gonic/gin"
)

func AddGeneralHandlers(r *gin.Engine) {
	// GET Routes
	r.GET("/", handlers.HomeHandler)
	r.GET("/network/interfaces", handlers.GetInterfacesHandler)
	r.GET("/network/interfaces/names", handlers.GetInterfaceNamesHandler)
	r.GET("/network/interface/:name", handlers.GetInterfaceByNameHandler)
	r.GET("/network/routes", handlers.GetRoutesHandler)
	r.GET("/secret", handlers.SecretTestHandler)

	// POST Routes
	r.POST("/login", handlers.LoginHandler)
}
