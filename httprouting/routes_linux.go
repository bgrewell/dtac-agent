// +build linux

package httprouting

import (
	"github.com/BGrewell/system-api/handlers"
	"github.com/gin-gonic/gin"
)

func AddOSSpecificHandlers(r *gin.Engine) {
	r.GET("/network/iptables/dnat", handlers.GetIptablesDnatRulesHandler)
	r.GET("examples/network/iptables/dnat/", handlers.GetIptablesDnatExamplesHandler)
	r.POST("/network/iptables/dnat/:id", handlers.CreateIptablesDnatRuleHandler)
	r.DELETE("/network/iptables/dnat/:id", handlers.DeleteIptablesDnatRuleHandler)
}
