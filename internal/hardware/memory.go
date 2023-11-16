package hardware

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

// MemoryInfo is the interface for the memory subsystem
type MemoryInfo interface {
	Update()
	Info() *mem.VirtualMemoryStat
}

// LiveMemoryInfo is the struct for the memory subsystem
type LiveMemoryInfo struct {
	Logger   *zap.Logger
	MemStats *mem.VirtualMemoryStat
}

// Update updates the memory subsystem
func (i *LiveMemoryInfo) Update() {
	n, err := mem.VirtualMemory()
	if err != nil {
		i.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	i.MemStats = n
}

// Info returns the memory subsystem info
func (i *LiveMemoryInfo) Info() *mem.VirtualMemoryStat {
	return i.MemStats
}

func (s *Subsystem) memInfoHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		s.mem.Update()
		return s.mem.Info(), nil
	}, "memory information")
}
