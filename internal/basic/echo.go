package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// NewEchoSubsystem creates a new echo subsystem
func NewEchoSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "echo"
	es := EchoSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    false, // should be added to configuration, hardcoded to false to disable
		name:       name,
	}
	return &es
}

// EchoSubsystem is a simple example subsystem for showing how the pieces fit together
type EchoSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	enabled    bool        // Optional subsystems have a boolean to control if they are enabled
	name       string      // Subsystem name
}

// Register registers the routes that this module handles
func (es *EchoSubsystem) Register() error {
	if !es.Enabled() {
		es.Logger.Info("subsystem is disabled", zap.String("subsystem", es.Name()))
		return nil
	}

	// Create a group for this subsystem
	base := es.Controller.Router.Group(es.name)

	// Routes
	secure := es.Controller.Config.Auth.DefaultSecure
	routes := []types.RouteInfo{
		{Group: base, HTTPMethod: http.MethodGet, Path: "/", Handler: es.rootHandler, Protected: secure},
	}

	// Register routes
	helpers.RegisterRoutes(routes, es.Controller.SecureMiddleware)
	es.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

// Enabled returns whether or not the echo subsystem is enabled
func (es *EchoSubsystem) Enabled() bool {
	return es.enabled
}

// Name returns the name of the echo subsystem
func (es *EchoSubsystem) Name() string {
	return es.name
}

func (es *EchoSubsystem) rootHandler(c *gin.Context) {
	start := time.Now()
	msg := "you must pass a name using ?name=<name>"
	if name := c.Query("name"); name != "" {
		msg = fmt.Sprintf("hello and welcome %s!", name)
	}
	response := gin.H{
		"message": msg,
	}
	es.Controller.Formatter.WriteResponse(c, time.Since(start), response)

}
