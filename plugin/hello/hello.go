package main

import (
	"github.com/BGrewell/system-api/plugin"
	"github.com/gin-gonic/gin"
	"time"
)

func Load() plugin.Plugin {
	return &HelloPlugin{}
}

type HelloPlugin struct {
	Message string
}

func (p *HelloPlugin) Register(config map[string]interface{}, r *gin.RouterGroup) error {
	// Load configuration if it exists
	p.Message = "hello default"
	if msg, ok := config["message"]; ok {
		p.Message = msg.(string)
	}

	// Register routes
	r.GET("/", p.HelloHandler)
	return nil
}

func (p *HelloPlugin) Name() string {
	return "hello"
}

func (p *HelloPlugin) HelloHandler(c *gin.Context) {
	start := time.Now()
	c.JSON(200, gin.H{
		"message": p.Message,
		"time":    start.Format("2006-01-02 15:04:05"),
	})
}
