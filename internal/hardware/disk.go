package hardware

import (
	"encoding/json"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/shirou/gopsutil/disk"
	"go.uber.org/zap"
)

// DiskUsageArgs is a struct to assist with validating the input arguments for disk usage
type DiskUsageArgs struct {
	Path string `json:"path,omitempty" yaml:"path,omitempty" xml:"path,omitempty"`
}

// DiskDetails is the struct for the disk details
type DiskDetails struct {
	Name   string
	Size   string
	Model  string
	Serial string
	Label  string
}

// DiskReport is the struct for the disk report
type DiskReport struct {
	Disks      []*DiskDetails       `json:"disks"`
	Partitions []disk.PartitionStat `json:"partitions"`
	Usage      []*disk.UsageStat    `json:"usage"`
}

// DiskInfo is the interface for the disk subsystem
type DiskInfo interface {
	Update()
	Info() *DiskReport
}

// LiveDiskInfo is the struct for the disk subsystem
type LiveDiskInfo struct {
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	DiskReport *DiskReport
}

// Update updates the disk subsystem
func (ni *LiveDiskInfo) Update() {
	p, err := disk.Partitions(true)
	if err != nil {
		ni.Logger.Error("failed to get disk stats", zap.Error(err))
	}

	d, err := GetPhysicalDisks()
	if err != nil {
		ni.Logger.Error("failed to get physical disks", zap.Error(err))
	}

	du := make([]*disk.UsageStat, 0)
	for _, partition := range p {
		u, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			ni.Logger.Error("failed to get disk usage", zap.Error(err))
		}
		du = append(du, u)
	}

	dr := DiskReport{
		Disks:      d,
		Partitions: p,
		Usage:      du,
	}
	ni.DiskReport = &dr
}

// Info returns the disk subsystem info
func (ni *LiveDiskInfo) Info() *DiskReport {
	return ni.DiskReport
}

func (s *Subsystem) diskRootHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		s.disk.Update()
		return json.Marshal(s.disk.Info())
	}, "disk information")
}

func (s *Subsystem) diskPartitionHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		s.disk.Update()
		return json.Marshal(s.disk.Info().Partitions)
	}, "disk partition information")
}

func (s *Subsystem) diskPhysicalDisksHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		s.disk.Update()
		return json.Marshal(s.disk.Info().Disks)
	}, "physical disk information")
}

func (s *Subsystem) diskUsageHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		path := ""
		if v, ok := in.Parameters["path"]; ok {
			path = v[0]
		}
		du := make([]*disk.UsageStat, 0)
		s.disk.Update()
		if path == "" {
			du = append(du, s.disk.Info().Usage...)
		} else {
			for _, stats := range s.disk.Info().Usage {
				if stats.Path == path {
					du = append(du, stats)
				}
			}
		}
		return json.Marshal(du)

	}, "disk usage information")
}
