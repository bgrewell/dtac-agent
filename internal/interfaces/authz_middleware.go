package interfaces

import "github.com/gin-gonic/gin"

// AuthorizationMiddleware is the interface for the authorization middleware
type AuthorizationMiddleware interface {
	AuthorizationHandler(*gin.Context)
}
