package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"time"
)

func getRoutesHandler(c *gin.Context) {
	start := time.Now()
	routes, err := GetRouteTable()
	if err != nil {
		helpers.WriteErrorResponseJSON(c, err)
		return
	}
	helpers.WriteResponseJSON(c, time.Since(start), routes)
}

func getRouteHandler(c *gin.Context) {
	helpers.WriteNotImplementedResponseJSON(c)
}

func createRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.WriteErrorResponseJSON(c, err)
		return
	}
	if err := CreateRoute(*input); err != nil {
		helpers.WriteErrorResponseJSON(c, fmt.Errorf("failed to create route: %v", err))
		return
	}
	output, err := GetRouteTable()
	if err != nil {
		helpers.WriteErrorResponseJSON(c, fmt.Errorf("route may not have been created. failed to retreive route table: %v", err))
		return
	}
	helpers.WriteResponseJSON(c, time.Since(start), output)
}

func updateRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.WriteErrorResponseJSON(c, err)
		return
	}
	if err := UpdateRoute(*input); err != nil {
		helpers.WriteErrorResponseJSON(c, fmt.Errorf("failed to update route: %v", err))
		return
	}
	output, err := GetRouteTable()
	if err != nil {
		helpers.WriteErrorResponseJSON(c, fmt.Errorf("route may not have been updated. failed to retreive route table: %v", err))
		return
	}
	helpers.WriteResponseJSON(c, time.Since(start), output)
}

func deleteRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.WriteErrorResponseJSON(c, err)
		return
	}
	if err := DeleteRoute(*input); err != nil {
		helpers.WriteErrorResponseJSON(c, fmt.Errorf("failed to delete route: %v", err))
		return
	}
	output, err := GetRouteTable()
	if err != nil {
		helpers.WriteErrorResponseJSON(c, fmt.Errorf("route may not have been deleted. failed to retreive route table: %v", err))
	}
	helpers.WriteResponseJSON(c, time.Since(start), output)
}

func getRoutesUnifiedHandler(c *gin.Context) {
	helpers.WriteNotImplementedResponseJSON(c)
}

func getRouteUnifiedHandler(c *gin.Context) {
	helpers.WriteNotImplementedResponseJSON(c)
}

func createRouteUnifiedHandler(c *gin.Context) {
	helpers.WriteNotImplementedResponseJSON(c)
}

func updateRouteUnifiedHandler(c *gin.Context) {
	helpers.WriteNotImplementedResponseJSON(c)
}

func deleteRouteUnifiedHandler(c *gin.Context) {
	helpers.WriteNotImplementedResponseJSON(c)
}
