package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authndb"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/endpoints"
	"go.uber.org/zap"
)

// Controller is the struct for the controller
type Controller struct {
	Logger           *zap.Logger
	Config           *config.Configuration
	EndpointList     *endpoints.EndpointList
	SecureMiddleware []gin.HandlerFunc
	AuthDB           *authndb.AuthDB
}
