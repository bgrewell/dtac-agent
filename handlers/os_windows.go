package handlers

import (
	"fmt"
	. "github.com/intel-innersource/frameworks.automation.dtac.agent/common"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/network"
	"github.com/BGrewell/go-execute"
	"github.com/BGrewell/go-netqospolicy"
	"github.com/gin-gonic/gin"
	"time"
)

func GetNetQosPoliciesHandler(c *gin.Context) {
	start := time.Now()
	policies, err := network.GetNetQosPolicies()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), policies)
}

func GetNetQosPolicyHandler(c *gin.Context) {
	start := time.Now()
	name := c.Param("name")
	if name != "" {
		policy, err := network.GetNetQosPolicy(name)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), policy)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func UpdateNetQosPolicyHandler(c *gin.Context) {
	start := time.Now()
	var input *netqospolicy.NetQoSPolicy
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if err := network.UpdateNetQosPolicy(input); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to update policy: %v", err))
		return
	}
	output, err := network.GetNetQosPolicy(input.Name)
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("policy may not have been updated. failed to retreive newly created policy: %v", err))
		return
	}
	WriteResponseJSON(c, time.Since(start), output)
}

func CreateNetQosPolicyHandler(c *gin.Context) {
	start := time.Now()
	var input *netqospolicy.NetQoSPolicy
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if err := network.CreateNetQosPolicy(input); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to create policy: %v", err))
		return
	}
	output, err := network.GetNetQosPolicy(input.Name)
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("policy may not have been created. failed to retreive newly created policy: %v", err))
	}
	WriteResponseJSON(c, time.Since(start), output)
}

func DeleteNetQosPoliciesHandler(c *gin.Context) {
	start := time.Now()
	var policies []*netqospolicy.NetQoSPolicy
	policies, _ = network.GetNetQosPolicies()
	if err := network.DeleteNetQosPolicies(); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to delete policies: %v", err))
		return
	}
	WriteResponseJSON(c, time.Since(start), policies)
}

func DeleteNetQosPolicyHandler(c *gin.Context) {
	start := time.Now()
	name := c.Param("name")
	if name != "" {
		var policy *netqospolicy.NetQoSPolicy
		policy, _ = network.GetNetQosPolicy(name)
		if err := network.DeleteNetQosPolicy(name); err != nil {
			WriteErrorResponseJSON(c, fmt.Errorf("failed to delete policy: %v", err))
			return
		}
		WriteResponseJSON(c, time.Since(start), policy)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func SystemRebootHandler(c *gin.Context) {
	start := time.Now()
	out, err := execute.ExecuteCmd("shutdown /r")
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
	WriteNotImplementedResponseJSON(c)
}
