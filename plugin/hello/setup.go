package hello

import "github.com/gin-gonic/gin"

type HelloPlugin struct {
	Message string
}

func (p *HelloPlugin) Register(config map[string]interface{}, r *gin.RouterGroup) error {
	p.Message = "hello default"
	if msg, ok := config["message"]; ok {
		p.Message = msg.(string)
	}
	r.GET("/", p.HelloHandler)
	return nil
}

func (p *HelloPlugin) Name() string {
	return "hello"
}
