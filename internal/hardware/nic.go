package hardware

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
)

func NewNicInfo(router *gin.Engine, cfg *config.Configuration, log *zap.Logger) *NicInfo {
	ni := NicInfo{}
	return &ni
}

type NicInfo struct {
}
