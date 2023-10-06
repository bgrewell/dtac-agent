package types

import "github.com/gin-gonic/gin"

type RouteInfo struct {
	Group      *gin.RouterGroup
	HttpMethod string
	Path       string
	Handler    gin.HandlerFunc
	Protected  bool
}
