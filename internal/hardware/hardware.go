package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
)

func NewHardwareInfo(router *gin.Engine, cfg *config.Configuration, log *zap.Logger) *HardwareInfo {
	hwi := HardwareInfo{}
	return &hwi
}

type HardwareInfo struct {
}
