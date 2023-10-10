package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
	"time"
)

type MemoryInfo interface {
	Update()
	Info() *mem.VirtualMemoryStat
}

type LiveMemoryInfo struct {
	Logger   *zap.Logger
	MemStats *mem.VirtualMemoryStat
}

func (i *LiveMemoryInfo) Update() {
	n, err := mem.VirtualMemory()
	if err != nil {
		i.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	i.MemStats = n
}

func (i *LiveMemoryInfo) Info() *mem.VirtualMemoryStat {
	return i.MemStats
}

func (s *HardwareSubsystem) memInfoHandler(c *gin.Context) {
	start := time.Now()
	s.mem.Update()
	helpers.WriteResponseJSON(c, time.Since(start), s.mem.Info())
}
