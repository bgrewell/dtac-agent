package hello

import "github.com/gin-gonic/gin"

type HelloPlugin struct {
	Message string
}

func (p *HelloPlugin) Register(r *gin.Engine) {
	// TODO: Should figure out how to make all routes have to be /<plugin_name>/...
	r.GET("/hello")
}
