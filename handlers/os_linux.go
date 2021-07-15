package handlers

import (
	"errors"
	"fmt"
	"github.com/BGrewell/go-execute"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/network"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func CreateIptablesDSCPRuleHandler(c *gin.Context) {
	start := time.Now()
	var input *network.DSCPRule
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
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
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
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func GetIptablesDSCPRuleHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func GetIptablesStatusHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func GetIptablesRulesHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
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

func GetIptablesRulesByChainHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func GetIptablesRuleByChainHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func CreateDNATRuleHandler(c *gin.Context) {
	start := time.Now()
	var input *network.DNATRule
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
		var input *network.DNATRule
		if err := c.ShouldBindJSON(&input); err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		rule, err := network.IptablesUpdateDNatRule(id, input)
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
	err := network.IptablesDelDNatRules()
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
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func GetDNATRuleHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func CreateSNATRuleHandler(c *gin.Context) {
	start := time.Now()
	var input *network.SNATRule
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
		var input *network.SNATRule
		if err := c.ShouldBindJSON(&input); err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		rule, err := network.IptablesUpdateSNatRule(id, input)
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
	err := network.IptablesDelSNatRules()
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
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
}

func GetSNATRuleHandler(c *gin.Context) {
	WriteErrorResponseJSON(c, errors.New("this method has not been implemented for this operating system yet"))
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