package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
)

func NewNicInfo(router *gin.Engine, cfg *config.Configuration, log *zap.Logger) NicInfo {
	ni := LiveNicInfo{
		Router: router,
		Config: cfg,
		Logger: log.With(zap.String("module", "network_info")),
	}
	n, err := net.Interfaces()
	if err != nil {
		ni.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	ni.InterfaceStats = n
	return &ni
}

type NicInfo interface {
	Update()
	Info() []net.InterfaceStat
}

type LiveNicInfo struct {
	Router         *gin.Engine           // All subsystems have a pointer to the gin.Engine
	Config         *config.Configuration // All subsystems have a pointer to the configuration
	Logger         *zap.Logger           // All subsystems have a pointer to the logger
	InterfaceStats []net.InterfaceStat
}

func (ni *LiveNicInfo) Update() {
	n, err := net.Interfaces()
	if err != nil {
		ni.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	ni.InterfaceStats = n
}

func (ni *LiveNicInfo) Info() []net.InterfaceStat {
	return ni.InterfaceStats
}
