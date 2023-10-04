package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
)

func NewEchoSubsystem(router *gin.Engine, log *zap.Logger, cfg *config.Configuration) *EchoSubsystem {
	es := EchoSubsystem{
		Router:  router,
		Logger:  log.With(zap.String("module", "echo")),
		Config:  cfg,
		Enabled: cfg.Subsystems.Diag,
	}
	if es.Enabled {
		err := es.Register()
		if err != nil {
			es.Logger.Error("failed to initialize plugin subsystem", zap.Error(err))
			return nil
		}
	}
	return &es
}

// EchoSubsystem is a simple example subsystem for showing how the pieces fit together
type EchoSubsystem struct {
	Enabled bool                  // Optional subsystems have a boolean to control if they are enabled
	Router  *gin.Engine           // All subsystems have a pointer to the gin.Engine
	Config  *config.Configuration // All subsystems have a pointer to the configuration
	Logger  *zap.Logger           // All subsystems have a pointer to the logger
}

// Register() registers the routes that this module handles
func (es *EchoSubsystem) Register() error {
	// Create a group for this subsystem
	base := es.Router.Group("echo")

	// Routes
	routes := []types.RouteInfo{
		{HttpMethod: http.MethodGet, Path: "/", Handler: es.rootHandler},
	}

	// Register routes
	for _, route := range routes {
		base.Handle(route.HttpMethod, route.Path, route.Handler)
	}
	es.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

func (es *EchoSubsystem) rootHandler(c *gin.Context) {
	msg := "you must pass a name using ?name=<name>"
	if name := c.Query("name"); name != "" {
		msg = fmt.Sprintf("hello and welcome %s!", name)
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": msg,
	})
}
