package interfaces

import "github.com/gin-gonic/gin"

// AuthenticationMiddleware is the interface for the authentication middleware
type AuthenticationMiddleware interface {
	AuthenticationHandler(*gin.Context)
}
