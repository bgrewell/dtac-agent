package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authn_db"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/http"
	"go.uber.org/zap"
)

// Controller is the struct for the controller
type Controller struct {
	Router           *gin.Engine
	Logger           *zap.Logger
	Config           *config.Configuration
	HttpRouteList    *http.HttpRouteList
	SecureMiddleware []gin.HandlerFunc
	AuthDB           *authn_db.AuthDB
}
