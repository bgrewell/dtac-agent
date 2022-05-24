package hello

import (
	. "github.com/BGrewell/dtac-agent/common"
	"github.com/gin-gonic/gin"
	"time"
)

func (p *HelloModule) HelloHandler(c *gin.Context) {
	start := time.Now()
	WriteResponseJSON(c, time.Since(start), p.Message)
}
