package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
	"strings"
	"time"
)

// CPUInfo is the interface for the cpu subsystem
type CPUInfo interface {
	Update()
	Info() []cpu.InfoStat
	Percent(interval time.Duration, perCPU bool) ([]float64, error)
}

// LiveCPUInfo is the struct for the cpu subsystem
type LiveCPUInfo struct {
	Logger   *zap.Logger // All subsystems have a pointer to the logger
	CPUStats []cpu.InfoStat
}

// Update updates the cpu subsystem
func (i *LiveCPUInfo) Update() {
	n, err := cpu.Info()
	if err != nil {
		i.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	i.CPUStats = n
}

// Info returns the cpu subsystem info
func (i *LiveCPUInfo) Info() []cpu.InfoStat {
	return i.CPUStats
}

// Percent returns the cpu subsystem percent
func (i *LiveCPUInfo) Percent(interval time.Duration, perCPU bool) ([]float64, error) {
	return cpu.Percent(interval, perCPU)
}

func (s *Subsystem) cpuInfoHandler(c *gin.Context) {
	start := time.Now()
	s.cpu.Update()
	helpers.WriteResponseJSON(c, time.Since(start), s.cpu.Info())
}

func (s *Subsystem) cpuUsageHandler(c *gin.Context) {
	perCore := true
	perCoreStr := c.Param("per_core")
	if perCoreStr != "" && strings.ToLower(perCoreStr) == "false" {
		perCore = false
	}
	start := time.Now()
	stats, err := cpu.Percent(time.Millisecond*100, perCore)
	if err != nil {
		helpers.WriteErrorResponseJSON(c, err)
	}
	helpers.WriteResponseJSON(c, time.Since(start), gin.H{
		"cpu_usage": stats,
	})
}
