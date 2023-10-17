package types

import "github.com/gin-gonic/gin"

// RouteInfo is a struct helper for registering routes
type RouteInfo struct {
	Group      *gin.RouterGroup
	HttpMethod string
	Path       string
	Handler    gin.HandlerFunc
	Protected  bool
}
