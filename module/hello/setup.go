package hello

import "github.com/gin-gonic/gin"

type HelloModule struct {
	Message string
}

func (p *HelloModule) Register(config map[string]interface{}, r *gin.RouterGroup) error {
	p.Message = "hello default"
	if msg, ok := config["message"]; ok {
		p.Message = msg.(string)
	}
	r.GET("/", p.HelloHandler)
	return nil
}

func (p *HelloModule) Name() string {
	return "hello"
}
