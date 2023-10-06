package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/shirou/gopsutil/net"
	"go.uber.org/zap"
)

func NewMockNicInfo() *MockNicInfo {
	mockData := []net.InterfaceStat{
		{
			MTU:          1500,
			Name:         "eth0",
			HardwareAddr: "00:1a:2b:3c:4d:5e",
			Flags:        []string{"up", "broadcast", "multicast"},
			Addrs: []net.InterfaceAddr{
				{
					Addr: "192.168.1.10/24",
				},
				{
					Addr: "fe80::21a:2bff:fe3c:4d5e/64",
				},
			},
		},
		{
			MTU:          1500,
			Name:         "eth1",
			HardwareAddr: "00:1a:2b:3c:4d:5f",
			Flags:        []string{"up", "broadcast", "multicast"},
			Addrs: []net.InterfaceAddr{
				{
					Addr: "192.168.1.11/24",
				},
			},
		},
	}

	return &MockNicInfo{
		InterfaceStats: mockData,
	}
}

type MockNicInfo struct {
	Router         *gin.Engine           // All subsystems have a pointer to the gin.Engine
	Config         *config.Configuration // All subsystems have a pointer to the configuration
	Logger         *zap.Logger           // All subsystems have a pointer to the logger
	InterfaceStats []net.InterfaceStat
}

func (ni *MockNicInfo) Update() {
	n, err := net.Interfaces()
	if err != nil {
		ni.Logger.Error("failed to get interface stats", zap.Error(err))
	}
	ni.InterfaceStats = n
}

func (ni *MockNicInfo) Info() []net.InterfaceStat {
	return ni.InterfaceStats
}
