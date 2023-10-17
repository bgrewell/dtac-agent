package hardware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
	"time"
)

// NicInfo is the interface for the nic subsystem
type NicInfo interface {
	Update()
	Info() []net.InterfaceStat
}

// LiveNicInfo is the struct for the nic subsystem
type LiveNicInfo struct {
	Logger         *zap.Logger // All subsystems have a pointer to the logger
	InterfaceStats []net.InterfaceStat
}

// Update updates the nic subsystem
func (ni *LiveNicInfo) Update() {
	n, err := net.Interfaces()
	if err != nil {
		ni.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	ni.InterfaceStats = n
}

// Info returns the nic subsystem info
func (ni *LiveNicInfo) Info() []net.InterfaceStat {
	return ni.InterfaceStats
}

// rootHandler handles requests for the root path for this subsystem
func (s *HardwareSubsystem) nicRootHandler(c *gin.Context) {
	start := time.Now()
	s.nic.Update()
	helpers.WriteResponseJSON(c, time.Since(start), s.nic.Info())
}

// nicInterfaceHandler handles requests for the root path for this subsystem
func (s *HardwareSubsystem) nicInterfaceHandler(c *gin.Context) {
	start := time.Now()
	name := c.Param("name")
	if name == "" {
		helpers.WriteErrorResponseJSON(c, errors.New("required path parameter 'name' not found. Ex: .../interface/<name>"))
		return
	}
	s.nic.Update()
	for _, info := range s.nic.Info() {
		if info.Name == name {
			helpers.WriteResponseJSON(c, time.Since(start), info)
			return
		}
	}
	helpers.WriteErrorResponseJSON(c, fmt.Errorf("no interface found by name: %s", name))
}
