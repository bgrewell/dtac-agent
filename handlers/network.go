package handlers

import (
	"fmt"
	. "github.com/intel-innersource/frameworks.automation.dtac.agent/common"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/network"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func GetInterfacesHandler(c *gin.Context) {
	start := time.Now()
	ifaces, err := network.GetInterfaces()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), ifaces)
}

func GetInterfaceNamesHandler(c *gin.Context) {
	start := time.Now()
	names, err := network.GetInterfaceNames()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), names)
}

func GetInterfaceByNameHandler(c *gin.Context) {
	start := time.Now()
	name := c.Param("name")
	if name != "" {
		var iface *network.Interface
		_, err := strconv.ParseInt(name, 10, 64)
		if err == nil {
			iface, err = network.GetInterfaceByIdx(name)
		} else {
			iface, err = network.GetInterfaceByName(name)
		}
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), iface)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
	}
}

func GetInterfaceByIdxHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		iface, err := network.GetInterfaceByIdx(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, time.Since(start), iface)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving id"))
	}
}

func GetRoutesHandler(c *gin.Context) {
	start := time.Now()
	routes, err := network.GetRouteTable()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), routes)
}

func CreateRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *network.RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if err := network.CreateRoute(*input); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to create route: %v", err))
		return
	}
	output, err := network.GetRouteTable()
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("route may not have been created. failed to retreive route table: %v", err))
		return
	}
	WriteResponseJSON(c, time.Since(start), output)
}

func UpdateRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *network.RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if err := network.UpdateRoute(*input); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to update route: %v", err))
		return
	}
	output, err := network.GetRouteTable()
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("route may not have been updated. failed to retreive route table: %v", err))
		return
	}
	WriteResponseJSON(c, time.Since(start), output)
}

func DeleteRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *network.RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	if err := network.DeleteRoute(*input); err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to delete route: %v", err))
		return
	}
	output, err := network.GetRouteTable()
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("route may not have been deleted. failed to retreive route table: %v", err))
	}
	WriteResponseJSON(c, time.Since(start), output)
}

func CreateFirewallRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func UpdateFirewallRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func DeleteFirewallRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func GetFirewallRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func GetFirewallRulesUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func CreateQosRuleUniversalHandler(c *gin.Context) {
	start := time.Now()
	var input *network.UniversalDSCPRule
	if err := c.ShouldBindJSON(&input); err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	output, err := network.CreateUniversalQosRule(input)
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to create universal dscp rule: %v", err))
		return
	}
	WriteResponseJSON(c, time.Since(start), output)
}

func GetQosRuleUniversalHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		output, err := network.GetUniversalQosRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, fmt.Errorf("failed to get universal dscp rule: %v", err))
			return
		}
		WriteResponseJSON(c, time.Since(start), output)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("required field 'id' not found"))
	}
}

func GetQosRulesUniversalHandler(c *gin.Context) {
	start := time.Now()
	output, err := network.GetUniversalQosRules()
	if err != nil {
		WriteErrorResponseJSON(c, fmt.Errorf("failed to get universal dscp rules: %v", err))
		return
	}
	WriteResponseJSON(c, time.Since(start), output)
}

func UpdateQosRuleUniversalHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		var input *network.UniversalDSCPRule
		if err := c.ShouldBindJSON(&input); err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		output, err := network.UpdateUniversalQosRule(id, input)
		if err != nil {
			WriteErrorResponseJSON(c, fmt.Errorf("failed to update universal dscp rule: %v", err))
			return
		}
		WriteResponseJSON(c, time.Since(start), output)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("required field 'id' not found"))
	}
}

func DeleteQosRuleUniversalHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id != "" {
		err := network.DeleteUniversalQosRule(id)
		if err != nil {
			WriteErrorResponseJSON(c, fmt.Errorf("failed to delete universal dscp rule: %v", err))
			return
		}
		WriteResponseJSON(c, time.Since(start), "{\"result\": \"deleted\"}")
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("required field 'id' not found"))
	}
}

func CreateRouteRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func GetRouteRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func GetRouteRulesUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func UpdateRouteRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func DeleteRouteRuleUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func GetInterfaceUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func GetInterfacesUniversalHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}
