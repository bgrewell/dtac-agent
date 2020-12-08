package handlers

import (
	"fmt"
	"github.com/BGrewell/go-netqospolicy"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/network"
	"github.com/gin-gonic/gin"
)


func GetNetQosPolicies(c *gin.Context) {
	policies, err := network.GetNetQosPolicies()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, policies)
}

func GetNetQosPolicy(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		policy, err := network.GetNetQosPolicy(name)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, policy)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}

func UpdateNetQosPolicy(c *gin.Context) {
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
	WriteResponseJSON(c, output)
}

func CreateNetQosPolicy(c *gin.Context) {
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
	WriteResponseJSON(c, output)
}

func DeleteNetQosPolicies(c *gin.Context) {
	var policies []*netqospolicy.NetQoSPolicy
	policies, _ = network.GetNetQosPolicies()
	if err := network.DeleteNetQosPolicies(); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to delete policies: %v", err))
		return
	}
	WriteResponseJSON(c, policies)
}

func DeleteNetQosPolicy(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		var policy *netqospolicy.NetQoSPolicy
		policy, _ = network.GetNetQosPolicy(name)
		if err := network.DeleteNetQosPolicy(name); err != nil {
			WriteErrorResponseJSON(c, fmt.Errorf("failed to delete policy: %v", err))
			return
		}
		WriteResponseJSON(c, policy)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
		return
	}
}