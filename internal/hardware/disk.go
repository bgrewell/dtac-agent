package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/shirou/gopsutil/disk"
	"go.uber.org/zap"
	"time"
)

type DiskDetails struct {
	Name   string
	Size   string
	Model  string
	Serial string
	Label  string
}

type DiskReport struct {
	Disks      []*DiskDetails       `json:"disks"`
	Partitions []disk.PartitionStat `json:"partitions"`
	Usage      []*disk.UsageStat    `json:"usage"`
}

type DiskInfo interface {
	Update()
	Info() *DiskReport
}

type LiveDiskInfo struct {
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	DiskReport *DiskReport
}

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

func (ni *LiveDiskInfo) Info() *DiskReport {
	return ni.DiskReport
}

// rootHandler handles requests for the root path for this subsystem
func (s *HardwareSubsystem) diskRootHandler(c *gin.Context) {
	start := time.Now()
	s.disk.Update()
	helpers.WriteResponseJSON(c, time.Since(start), s.disk.Info())
}

func (s *HardwareSubsystem) diskPartitionHandler(c *gin.Context) {
	start := time.Now()
	s.disk.Update()
	helpers.WriteResponseJSON(c, time.Since(start), s.disk.Info().Partitions)
}

func (s *HardwareSubsystem) diskPhysicalDisksHandler(c *gin.Context) {
	start := time.Now()
	s.disk.Update()
	helpers.WriteResponseJSON(c, time.Since(start), s.disk.Info().Disks)
}

func (s *HardwareSubsystem) diskUsageHandler(c *gin.Context) {
	start := time.Now()
	path := c.Query("path")
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
	helpers.WriteResponseJSON(c, time.Since(start), du)
}
