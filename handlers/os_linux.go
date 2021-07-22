package handlers

import (
	"fmt"
	"github.com/BGrewell/go-execute"
	"github.com/BGrewell/go-iptables"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/network"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func CreateIptablesDSCPRuleHandler(c *gin.Context) {
	start := time.Now()
	var input *network.DSCPTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	id, err := network.IptablesAddDSCPRule(input)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), id)
}

func UpdateIptablesDSCPRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	var input *network.DSCPTemplate
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if id != "" {
		rule, err := network.IptablesUpdateDSCPRule(id, input)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func DeleteIptablesDSCPRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesDelDSCPRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func GetIptablesDSCPRulesHandler(c *gin.Context) {
	start := time.Now()
	rules, err := network.IptablesGetDSCPRules()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), rules)
}

func GetIptablesDSCPRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesGetDSCPRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func GetIptablesStatusHandler(c *gin.Context) {
	start := time.Now()
	status, err := network.IptablesGetStatus()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), status)
}

func GetIptablesRulesHandler(c *gin.Context) {
	start := time.Now()
	rules, err := network.IptablesGetRules()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), rules)
}

func GetIptablesRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesGetRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func GetIptablesRulesByTableHandler(c *gin.Context) {
	start := time.Now()
	table := c.Param("table")
	if table != "" {
		rule, err := network.IptablesGetByTable(table)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func GetIptablesRulesByChainHandler(c *gin.Context) {
	start := time.Now()
	table := c.Param("table")
	chain := c.Param("chain")
	if table != "" && chain != "" {
		rule, err := network.IptablesGetByChain(table, chain)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error table and chain are required parameters"))
		return
	}
}

func DeleteIptablesRulesHandler(c *gin.Context) {
	start := time.Now()
	err := network.IptablesDelRules()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), "{\"result\": \"all application specific rules were deleted\"}")
}

func DeleteIptablesRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesDelRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error id is a required parameter"))
		return
	}
}

func CreateIptablesRuleHandler(c *gin.Context) {
	start := time.Now()
	var input *iptables.Rule
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}

	rule, err := network.IptablesCreateRule(input)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), rule)

}

func UpdateIptablesRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	var input *iptables.Rule
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if id != "" {
		rule, err := network.IptablesUpdateRule(id, input)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving id"))
		return
	}
}

func CreateDNATRuleHandler(c *gin.Context) {
	start := time.Now()
	var input *iptables.Rule
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	id, err := network.IptablesAddDNatRule(input)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), id)
}

func UpdateDNATRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		var input *iptables.Rule
		if err := c.ShouldBindJSON(&input); err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		rule, err := network.IptablesUpdateDNatRule(input)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func DeleteDNATRulesHandler(c *gin.Context) {
	start := time.Now()
	_, err := network.IptablesDelDNatRules()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), "{status: ok}")
}

func DeleteDNATRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesDelDNatRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func GetDNATRulesHandler(c *gin.Context) {
	start := time.Now()
	rules, err := network.IptablesGetDNatRules()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), rules)
}

func GetDNATRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesGetDNatRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func CreateSNATRuleHandler(c *gin.Context) {
	start := time.Now()
	var input *iptables.Rule
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	id, err := network.IptablesAddSNatRule(input)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), id)
}

func UpdateSNATRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		var input *iptables.Rule
		if err := c.ShouldBindJSON(&input); err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		rule, err := network.IptablesUpdateSNatRule(input)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func DeleteSNATRulesHandler(c *gin.Context) {
	start := time.Now()
	_, err := network.IptablesDelSNatRules()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), "{status: ok}")
}

func DeleteSNATRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesDelSNatRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func GetSNATRulesHandler(c *gin.Context) {
	start := time.Now()
	rules, err := network.IptablesGetSNatRules()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), rules)
}

func GetSNATRuleHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		rule, err := network.IptablesGetSNatRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), rule)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving id"))
		return
	}
}

func SystemRebootHandler(c *gin.Context) {
	start := time.Now()
	out, err := execute.ExecuteCmd("shutdown -r")
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to reboot computer: %v"))
		return
	}
	WriteResponseJSON(c, time.Since(start), out)
}

func SystemShutdownHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func SystemApiRestartHandler(c *gin.Context) {
	ts := c.Param("time")
	t := 10
	if ts != "" {
		var err error
		t, err = strconv.Atoi(ts)
		if err != nil {
			t = 10
		}
	}
	start := time.Now()
	go func() {
		time.Sleep(time.Duration(t) * time.Second)
		execute.ExecuteCmd("/bin/systemctl restart system-apid")
	}()
	WriteResponseJSON(c, time.Since(start), fmt.Sprintf("service will restart in %d seconds", t))
}