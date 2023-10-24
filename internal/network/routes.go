package network

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func (s *Subsystem) getRoutesHandler(c *gin.Context) {
	start := time.Now()
	routes, err := GetRouteTable()
	if err != nil {
		s.Controller.Formatter.WriteError(c, err)
		return
	}
	s.Controller.Formatter.WriteResponse(c, time.Since(start), routes)
}

func (s *Subsystem) getRouteHandler(c *gin.Context) {
	s.Controller.Formatter.WriteNotImplementedError(c, errors.New("this function has not been implemented yet"))
}

func (s *Subsystem) createRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		s.Controller.Formatter.WriteError(c, err)
		return
	}
	if err := CreateRoute(*input); err != nil {
		s.Controller.Formatter.WriteError(c, fmt.Errorf("failed to create route: %v", err))
		return
	}
	output, err := GetRouteTable()
	if err != nil {
		s.Controller.Formatter.WriteError(c, fmt.Errorf("route may not have been created. failed to retreive route table: %v", err))
		return
	}
	s.Controller.Formatter.WriteResponse(c, time.Since(start), output)
}

func (s *Subsystem) updateRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		s.Controller.Formatter.WriteError(c, err)
		return
	}
	if err := UpdateRoute(*input); err != nil {
		s.Controller.Formatter.WriteError(c, fmt.Errorf("failed to update route: %v", err))
		return
	}
	output, err := GetRouteTable()
	if err != nil {
		s.Controller.Formatter.WriteError(c, fmt.Errorf("route may not have been updated. failed to retreive route table: %v", err))
		return
	}
	s.Controller.Formatter.WriteResponse(c, time.Since(start), output)
}

func (s *Subsystem) deleteRouteHandler(c *gin.Context) {
	start := time.Now()
	var input *RouteTableRow
	if err := c.ShouldBindJSON(&input); err != nil {
		s.Controller.Formatter.WriteError(c, err)
		return
	}
	if err := DeleteRoute(*input); err != nil {
		s.Controller.Formatter.WriteError(c, fmt.Errorf("failed to delete route: %v", err))
		return
	}
	output, err := GetRouteTable()
	if err != nil {
		s.Controller.Formatter.WriteError(c, fmt.Errorf("route may not have been deleted. failed to retreive route table: %v", err))
	}
	s.Controller.Formatter.WriteResponse(c, time.Since(start), output)
}

func (s *Subsystem) getRoutesUnifiedHandler(c *gin.Context) {
	s.Controller.Formatter.WriteNotImplementedError(c, errors.New("this function has not been implemented yet"))
}

func (s *Subsystem) getRouteUnifiedHandler(c *gin.Context) {
	s.Controller.Formatter.WriteNotImplementedError(c, errors.New("this function has not been implemented yet"))
}

func (s *Subsystem) createRouteUnifiedHandler(c *gin.Context) {
	s.Controller.Formatter.WriteNotImplementedError(c, errors.New("this function has not been implemented yet"))
}

func (s *Subsystem) updateRouteUnifiedHandler(c *gin.Context) {
	s.Controller.Formatter.WriteNotImplementedError(c, errors.New("this function has not been implemented yet"))
}

func (s *Subsystem) deleteRouteUnifiedHandler(c *gin.Context) {
	s.Controller.Formatter.WriteNotImplementedError(c, errors.New("this function has not been implemented yet"))
}
