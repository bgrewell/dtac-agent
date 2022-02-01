package hello

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/gin-gonic/gin"
	"time"
)

func (p *HelloPlugin) HelloHandler(c *gin.Context) {
	start := time.Now()
	WriteResponseJSON(c, time.Since(start), p.Message)
}
