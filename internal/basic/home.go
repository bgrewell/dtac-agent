package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version"
	"go.uber.org/zap"
	"time"
)

// NewHomePageSubsystem creates a new homepage subsystem
func NewHomePageSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "homepage"
	hps := HomePageSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		name:       name,
	}
	return &hps
}

// HomePageSubsystem serves the main homepage content
type HomePageSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	name       string      // Subsystem name
}

// Register registers the routes for the homepage
func (hps *HomePageSubsystem) Register() error {
	if !hps.Enabled() {
		hps.Logger.Info("subsystem is disabled", zap.String("subsystem", hps.Name()))
		return nil
	}

	// Registering a route for the homepage
	hps.Controller.Router.GET("/", hps.homeHandler)
	hps.Logger.Info("homepage route registered")
	return nil
}

// Enabled returns whether or not the homepage subsystem is enabled
func (hps *HomePageSubsystem) Enabled() bool {
	return true
}

// Name returns the name of the homepage subsystem
func (hps *HomePageSubsystem) Name() string {
	return hps.name
}

func (hps *HomePageSubsystem) homeHandler(c *gin.Context) {
	start := time.Now()
	hps.Controller.HTTPRouteList.UpdateRoutes()
	response := gin.H{
		"message": fmt.Sprintf("welcome to the %s", hps.Controller.Config.Internal.ProductName),
		"version": version.Current().String(),
		"routes":  hps.Controller.HTTPRouteList.Routes,
	}
	hps.Controller.Formatter.WriteResponse(c, time.Since(start), response)
}
