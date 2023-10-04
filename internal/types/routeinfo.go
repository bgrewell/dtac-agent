package types

import "github.com/gin-gonic/gin"

type RouteInfo struct {
	HttpMethod string
	Path       string
	Handler    gin.HandlerFunc
}
