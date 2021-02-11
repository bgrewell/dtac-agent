package handlers

import (
	"fmt"
	"github.com/BGrewell/iptables"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/network"
	"github.com/gin-gonic/gin"
	"time"
)

func GetIptablesDnatRulesHandler(c *gin.Context) {

}

func GetIptablesDnatExamplesHandler(c *gin.Context) {

}

func CreateIptablesDnatRuleHandler(c *gin.Context) {
	start := time.Now()
	var rule *iptables.Rule
	if err := c.ShouldBindJSON(rule); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if id, err := network.AddIptablesDNatRule(rule); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to apply rule: %v", err))
	} else {
		WriteResponseJSON(c, time.Since(start), id)
	}
}

func DeleteIptablesDnatRuleHandler(c *gin.Context) {

}