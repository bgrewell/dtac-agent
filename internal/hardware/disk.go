package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/disk"
	"go.uber.org/zap"
	"time"
)

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

// rootHandler handles requests for the root path for this subsystem
func (s *Subsystem) diskRootHandler(c *gin.Context) {
	start := time.Now()
	s.disk.Update()
	s.Controller.Formatter.WriteResponse(c, time.Since(start), s.disk.Info())
}

func (s *Subsystem) diskPartitionHandler(c *gin.Context) {
	start := time.Now()
	s.disk.Update()
	s.Controller.Formatter.WriteResponse(c, time.Since(start), s.disk.Info().Partitions)
}

func (s *Subsystem) diskPhysicalDisksHandler(c *gin.Context) {
	start := time.Now()
	s.disk.Update()
	s.Controller.Formatter.WriteResponse(c, time.Since(start), s.disk.Info().Disks)
}

func (s *Subsystem) diskUsageHandler(c *gin.Context) {
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
	s.Controller.Formatter.WriteResponse(c, time.Since(start), du)
}
