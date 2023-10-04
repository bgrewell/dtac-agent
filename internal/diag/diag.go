package diag

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"go.uber.org/zap"
)

// NewSubsystem creates a new instances of the DiagSubsystem and if that subsystem is enabled it calls
// the Register() function to register the routes that the DiagSubsystem handles
func NewSubsystem(router *gin.Engine, log *zap.Logger, cfg *config.Configuration) *DiagSubsystem {
	ds := DiagSubsystem{
		Router:  router,
		Logger:  log.With(zap.String("module", "diag")),
		Config:  cfg,
		enabled: cfg.Subsystems.Diag,
	}
	if ds.enabled {
		err := ds.Register()
		if err != nil {
			ds.Logger.Error("failed to initialize plugin subsystem", zap.Error(err))
			return nil
		}
	}
	return &ds
}

// DiagSubsystem is the subsystem that contains routes related to internal dtac diagnostics
type DiagSubsystem struct {
	Router  *gin.Engine
	Logger  *zap.Logger
	Config  *config.Configuration
	enabled bool
}

// Register() registers the routes that this module handles
func (ds *DiagSubsystem) Register() error {
	// Create a group for this subsystem
	base := ds.Router.Group("diag")

	// Routes
	routes := []types.RouteInfo{
		{HttpMethod: http.MethodGet, Path: "/", Handler: ds.rootHandler},
	}

	// Register routes
	for _, route := range routes {
		base.Handle(route.HttpMethod, route.Path, route.Handler)
	}
	ds.Logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
}

// Enabled returns true if this module is enabled otherwise it returns false
func (ds *DiagSubsystem) Enabled() bool {
	return ds.enabled
}

// rootHandler handles requests for the root path for this subsystem
func (ds *DiagSubsystem) rootHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"version": types.AnnotatedStruct{
			Description: "dtac version information",
			Value:       version.Current(),
		},
		"memory": types.AnnotatedStruct{
			Description: "current dtac agent memory usage",
			Value:       CurrentMemoryStats(),
		},
	})
}

func (ds *DiagSubsystem) bobbyHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "this works",
	})
}
