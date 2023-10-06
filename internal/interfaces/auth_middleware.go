package interfaces

import "github.com/gin-gonic/gin"

type AuthMiddleware interface {
	AuthHandler(*gin.Context)
}
