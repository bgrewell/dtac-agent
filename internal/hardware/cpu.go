package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
	"strings"
	"time"
)

// CpuInfo is the interface for the cpu subsystem
type CpuInfo interface {
	Update()
	Info() []cpu.InfoStat
	Percent(interval time.Duration, perCpu bool) ([]float64, error)
}

// LiveCpuInfo is the struct for the cpu subsystem
type LiveCpuInfo struct {
	Logger   *zap.Logger // All subsystems have a pointer to the logger
	CpuStats []cpu.InfoStat
}

// Update updates the cpu subsystem
func (i *LiveCpuInfo) Update() {
	n, err := cpu.Info()
	if err != nil {
		i.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	i.CpuStats = n
}

// Info returns the cpu subsystem info
func (i *LiveCpuInfo) Info() []cpu.InfoStat {
	return i.CpuStats
}

// Percent returns the cpu subsystem percent
func (i *LiveCpuInfo) Percent(interval time.Duration, perCpu bool) ([]float64, error) {
	return cpu.Percent(interval, perCpu)
}

func (s *HardwareSubsystem) cpuInfoHandler(c *gin.Context) {
	start := time.Now()
	s.cpu.Update()
	helpers.WriteResponseJSON(c, time.Since(start), s.cpu.Info())
}

func (s *HardwareSubsystem) cpuUsageHandler(c *gin.Context) {
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
