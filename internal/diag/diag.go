package diag

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/register"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
)

// NewSubsystem creates a new instances of the DiagSubsystem and if that subsystem is enabled it calls
// the Register() function to register the routes that the DiagSubsystem handles
func NewSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "diag"
	ds := DiagSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    c.Config.Subsystems.Diag,
		name:       name,
	}
	return &ds
}

// DiagSubsystem is the subsystem that contains routes related to internal dtac diagnostics
type DiagSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger
	enabled    bool
	name       string // Subsystem name
}

// Register() registers the routes that this module handles
func (ds *DiagSubsystem) Register() error {
	if !ds.Enabled() {
		ds.Logger.Info("subsystem is disabled", zap.String("subsystem", ds.Name()))
		return nil
	}
	// Create a group for this subsystem
	base := ds.Controller.Router.Group(ds.name)

	// Routes
	secure := ds.Controller.Config.Auth.DefaultSecure
	routes := []types.RouteInfo{
		{Group: base, HttpMethod: http.MethodGet, Path: "/", Handler: ds.rootHandler, Protected: secure},
		{Group: base, HttpMethod: http.MethodGet, Path: "/jwt", Handler: ds.jwtTestHandler, Protected: true},
		{Group: base, HttpMethod: http.MethodGet, Path: "/routes", Handler: ds.httpRoutePrintHandler, Protected: secure},
	}

	// Register routes
	register.RegisterRoutes(routes, ds.Controller.SecureMiddleware)
	ds.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

// Enabled returns true if this module is enabled otherwise it returns false
func (ds *DiagSubsystem) Enabled() bool {
	return ds.enabled
}

func (ds *DiagSubsystem) Name() string {
	return ds.name
}

// rootHandler handles requests for the root path for this subsystem
func (ds *DiagSubsystem) rootHandler(c *gin.Context) {
	start := time.Now()
	response := gin.H{
		"version": types.AnnotatedStruct{
			Description: fmt.Sprintf("%s version information", ds.Controller.Config.Internal.ShortName),
			Value:       version.Current(),
		},
		"memory": types.AnnotatedStruct{
			Description: fmt.Sprintf("current %s memory usage", ds.Controller.Config.Internal.ShortName),
			Value:       CurrentMemoryStats(),
		},
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}

func (ds *DiagSubsystem) httpRoutePrintHandler(c *gin.Context) {
	start := time.Now()
	ds.Controller.HttpRouteList.UpdateRoutes()
	response := gin.H{
		"routes": types.AnnotatedStruct{
			Description: fmt.Sprintf("list of registered http endpoints being served by %s", ds.Controller.Config.Internal.ShortName),
			Value:       ds.Controller.HttpRouteList.Routes,
		},
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}

func (ds *DiagSubsystem) jwtTestHandler(c *gin.Context) {
	start := time.Now()
	response := gin.H{
		"message": "jwt test page",
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}
