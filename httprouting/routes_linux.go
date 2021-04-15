// +build linux

package httprouting

import (
	"github.com/gin-gonic/gin"
)

func AddOSSpecificHandlers(r *gin.Engine) {
	// IPTables DSCP
	r.POST("/network/dscp", handlers.CreateIptablesDSCPRuleHandler)       // Create IPTables DSCP Marking rule
	r.PUT("/network/dscp/:id", handlers.UpdateIptablesDSCPRuleHandler)    // Update IPTables DSCP Marking rule
	r.DELETE("/network/dscp/:id", handlers.DeleteIptablesDSCPRuleHandler) // Delete IPTables DSCP Marking rule
	r.GET("/network/dscp", handlers.GetIptablesDSCPRulesHandler)          // Get IPTables DSCP Marking rules
	r.GET("/network/dscp/:id", handlers.GetIptablesDSCPRuleHandler)       // Get IPTables DSCP Marking rules

	// IPTables DNAT
	r.POST("/network/dnat", handlers.CreateDNATRuleHandler)       // Create IPTables DNAT Rule
	r.PUT("/network/dnat/:id", handlers.UpdateDNATRuleHandler)    // Update IPTables DNAT Rule
	r.DELETE("/network/dnat/:id", handlers.DeleteDNATRuleHandler) // Delete IPTables DNAT Rule
	r.GET("/network/dnat", handlers.GetDNATRulesHandler)          // Get All DNAT Rules
	r.GET("/network/dnat/:id", handlers.GetDNATRuleHandler)       // Get DNAT Rule specified by id

	// IPTables SNAT
	r.POST("/network/snat", handlers.CreateSNATRuleHandler)       // Create IPTables SNAT Rule
	r.PUT("/network/snat/:id", handlers.UpdateSNATRuleHandler)    // Update IPTables SNAT Rule
	r.DELETE("/network/snat/:id", handlers.DeleteSNATRuleHandler) // Delete IPTables SNAT Rule
	r.GET("/network/snat", handlers.GetSNATRulesHandler)          // Get All SNAT Rules
	r.GET("/network/snat/:id", handlers.GetSNATRuleHandler)       // Get SNAT Rule specified by id

	// IPTables General
	r.GET("/network/firewall", handlers.GetIptablesStatusHandler)                       // Get IPTables Status
	r.GET("/network/firewall/rules", handlers.GetIptablesRulesHandler)                  // Get IPTables Rules
	r.GET("/network/firewall/rules/:chain", handlers.GetIptablesRulesByChainHandler)    // Get IPTables Rules by Chain
	r.GET("/network/firewall/rules/:chain/:id", handlers.GetIptablesRuleByChainHandler) // Get IPTables Rule by Id
}
