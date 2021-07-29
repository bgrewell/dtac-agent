package handlers

import (
	"errors"
	. "github.com/BGrewell/system-api/common"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var (
	Routes     gin.RoutesInfo
	Info       BasicInfo
)

func init() {
	Info = BasicInfo{}
	Info.Update()
}

func SecretTestHandler(c *gin.Context) {
	user, err := AuthorizeUser(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user.ID, "secret": "somesupersecretvalue"})
}

func HomeHandler(c *gin.Context) {
	// Update Routes
	start := time.Now()
	Info.UpdateRoutes(Routes)
	WriteResponseJSON(c, time.Since(start), Info)
}

func UnimplementedHandler(c *gin.Context) {
	WriteNotImplementedResponseJSON(c)
}

func GetCPUHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.CPU)
}

func GetCPUUsageHandler(c *gin.Context) {
	type CpuUsage struct {
		Combined float64 `json:"combined"`
		PerCore []float64 `json:"per_core"`
	}

	start := time.Now()
	percpu, err := cpu.Percent(0, true)
	if err != nil {
		WriteErrorResponseJSON(c, err)
	}
	time.Sleep(10 * time.Millisecond)
	total, err := cpu.Percent(0, false)
	if err != nil {
		WriteErrorResponseJSON(c, err)
	}
	usage := CpuUsage{
		Combined: total[0],
		PerCore:  percpu,
	}
	WriteResponseJSON(c, time.Since(start), usage)
}

func GetMemoryHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.Memory)
}

func GetHostHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.Host)
}

func GetNetworkHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.Network)
}

func GetEndpointsHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.Routes)
}

func GetDiskHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.Disk)
}

func GetDockerHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.Docker)
}

func GetProcessesHandler(c *gin.Context) {
	start := time.Now()
	Info.Update()
	WriteResponseJSON(c, time.Since(start), Info.Processes)
}

func GetProcessHandler(c *gin.Context) {
	start := time.Now()
	pids := c.Param("pid")
	if pids == "" {
		WriteErrorResponseJSON(c, errors.New("missing required pid"))
		return
	}
	pid, err := strconv.Atoi(pids)
	if err != nil {
		WriteErrorResponseJSON(c, errors.New("failed to convert pid to a number"))
		return
	}
	p, err := process.NewProcess(int32(pid))
	pi := ProcessInfo{}
	pi.Update(p)
	WriteResponseJSON(c, time.Since(start), pi)
}

func GetOSHandler(c *gin.Context) {
	type OSType struct {
		OperatingSystem string `json:"operating_system"`
	}
	start := time.Now()
	ost := OSType{ OperatingSystem: runtime.GOOS }
	WriteResponseJSON(c, time.Since(start), ost)
}