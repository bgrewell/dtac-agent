package interfaces

import "github.com/gin-gonic/gin"

type AuthorizationMiddleware interface {
	AuthorizationHandler(*gin.Context)
}
