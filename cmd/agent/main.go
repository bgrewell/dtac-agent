package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/adapters/rest"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authn"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authndb"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authz"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/basic"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/diag"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/endpoints"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/hardware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/middleware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/network"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/system"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"os"
	"runtime"
)

// AdapterParams is a struct that is used to pass the subsystems to any API frontends
type AdapterParams struct {
	fx.In
	LC         fx.Lifecycle
	Controller *controller.Controller
	Adapters   []interfaces.APIAdapter `group:"adapters"`
	Subsystems []interfaces.Subsystem  `group:"subsystems"`
}

// TODO: How to implement authn/authz middleware in such a way that it works with any API

// NewController creates a new instance of the controller.Controller struct
func NewController(logger *zap.Logger, cfg *config.Configuration, epl *endpoints.EndpointList,
	db *authndb.AuthDB) *controller.Controller {
	// Create the SubsystemParams object
	c := controller.Controller{
		Logger:           logger,
		Config:           cfg,
		EndpointList:     epl,
		SecureMiddleware: make([]gin.HandlerFunc, 0),
		AuthDB:           db,
	}
	return &c
}

// Setup is a function that is used to set up the API adapters
func Setup(params AdapterParams) {
	// Find middleware
	params.Controller.Logger.Debug("searching for middleware")
	middlewares := []middleware.Middleware{}
	for _, subsystem := range params.Subsystems {
		if sub, ok := subsystem.(middleware.Middleware); ok {
			params.Controller.Logger.Debug("found middleware", zap.String("name", sub.Name()))
			middlewares = append(middlewares, sub)
		}
	}

	// Prioritize middlewares
	params.Controller.Logger.Debug("prioritizing middleware")
	middlewares = middleware.Sort(middlewares)

	// Register middleware
	params.Controller.Logger.Debug("registering middleware", zap.Int("count", len(middlewares)))
	for _, subsystem := range params.Subsystems {
		endpionts := subsystem.Endpoints()
		for _, endpoint := range endpionts {
			endpoint.Function = middleware.Chain(middlewares, *endpoint)
		}
	}

	params.Controller.Logger.Debug("setting up API adapters", zap.Int("count", len(params.Adapters)))
	for _, adapter := range params.Adapters {
		params.Controller.Logger.Debug("registering adapter")
		params.Controller.Logger.Debug("registering adapter", zap.String("name", adapter.Name()))
		err := adapter.Register(params.Subsystems)
		if err != nil {
			params.Controller.Logger.Fatal("failed to register adapter", zap.Error(err))
		}
	}

	// Set up the Fx lifecycle controller
	params.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for _, adapter := range params.Adapters {
				params.Controller.Logger.Debug("starting adapter", zap.String("name", adapter.Name()))
				err := adapter.Start(ctx)
				if err != nil {
					params.Controller.Logger.Fatal("failed to start adapter", zap.Error(err))
				}
			}
			return nil
		},

		OnStop: func(ctx context.Context) error {
			for _, adapter := range params.Adapters {
				params.Controller.Logger.Debug("stopping adapter", zap.String("name", adapter.Name()))
				err := adapter.Stop(ctx)
				if err != nil {
					params.Controller.Logger.Fatal("failed to stop adapter", zap.Error(err))
				}
			}
			return nil
		},
	})
}

// AsSubsystem is a helper function that is used to annotate a function as a subsystem
func AsSubsystem(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(interfaces.Subsystem)),
		fx.ResultTags(`group:"subsystems"`),
	)
}

// AsAdapter is a helper function that is used to annotate a function as an API adapter
func AsAdapter(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(interfaces.APIAdapter)),
		fx.ResultTags(`group:"adapters"`))
}

func main() {

	if !helpers.IsRunningAsRoot() {
		if runtime.GOOS == "windows" {
			fmt.Println("Please run this application with elevated privileges!")
		} else {
			fmt.Println("Please run this application as root!")
		}
		os.Exit(1)
	}

	fx.New(
		// Setup zap logger with the Fx framework
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		// Set up the providers
		fx.Provide(
			config.NewConfiguration,                 // Configuration
			zap.NewDevelopment,                      // Structured Logger
			basic.NewTLSInfo,                        // Tls Cert Handler
			helpers.NewJSONResponseFormatter,        // Response Formatter
			endpoints.NewEndpointList,               // Endpoint List
			NewController,                           // Wrapper around common subsystem input components
			authndb.NewAuthDB,                       // Authentication database
			AsAdapter(rest.NewAdapter),              // Rest API Interface
			AsSubsystem(basic.NewHomePageSubsystem), // Homepage handler
			AsSubsystem(basic.NewEchoSubsystem),     // Demo Subsystem
			AsSubsystem(diag.NewSubsystem),          // Diagnostic Subsystem
			AsSubsystem(authn.NewSubsystem),         // Authentication Subsystem
			AsSubsystem(authz.NewSubsystem),         // Authorization Subsystem
			//AsSubsystem(plugin.NewSubsystem),        // Plugin Subsystem
			AsSubsystem(network.NewSubsystem),  // Network Subsystem
			AsSubsystem(hardware.NewSubsystem), // Hardware Subsystem
			AsSubsystem(system.NewSubsystem),   // System Subsystem
		),
		// Invoke any functions needed to initialize everything. The empty anonymous functions are
		// used to ensure that the providers that return that type are initialized.
		fx.Invoke(
			//authorization.EnsureAuthzModel,  // Ensure we have at least a default authorization model
			//authorization.EnsureAuthzPolicy, // Ensure we have at least a default authorization policy
			Setup, // Set up the application
		),
	).Run()
}
