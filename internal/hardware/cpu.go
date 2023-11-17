package hardware

import (
	"encoding/json"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/zap"
	"strings"
	"time"
)

// CPUUsageArgs is a struct to assist with validating the input arguments
type CPUUsageArgs struct {
	PerCore string `json:"per_core,omitempty" yaml:"per_core,omitempty" xml:"per_core,omitempty"`
}

// CPUUsageOutput is a struct to assist with describing the output format
type CPUUsageOutput struct {
	Usage []float64 `json:"usage,omitempty" yaml:"usage,omitempty" xml:"usage,omitempty"`
}

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

func (s *Subsystem) cpuInfoHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		s.cpu.Update()
		return json.Marshal(s.cpu.Info())
	}, "cpu information")
}

func (s *Subsystem) cpuUsageHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		perCore := true
		if v, ok := in.Parameters["per_core"]; ok {
			if v[0] != "" && strings.ToLower(v[0]) == "false" {
				perCore = false
			}
		}
		usage, err := cpu.Percent(time.Millisecond*100, perCore)
		if err != nil {
			return nil, err
		}
		return json.Marshal(CPUUsageOutput{
			Usage: usage,
		})
	}, "cpu usage information")
}
