package main

import (
	"context"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authn"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authndb"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/authz"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/basic"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config/authorization"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/hardware"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	httpRoutes "github.com/intel-innersource/frameworks.automation.dtac.agent/internal/http"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/network"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/plugin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/system"
	"go.uber.org/fx"
	"net"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/diag"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// RegisterParams is a struct that is used to pass the controller and subsystems to the Register function
type RegisterParams struct {
	fx.In
	Controller *controller.Controller
	Subsystems []interfaces.Subsystem `group:"subsystems"`
}

// NewHTTPServer creates the webserver that handles the requests sent to the DTAC agent
func NewHTTPServer(lc fx.Lifecycle, router *gin.Engine, log *zap.Logger, tls *basic.TLSInfo) *http.Server {

	// Create a new http server
	srv := &http.Server{Addr: ":8180", Handler: router}

	// Set up the serve function
	srvFunc := srv.Serve
	srvMsg := "starting HTTP server"

	if tls.Enabled {
		wrapper := func(l net.Listener) error {
			return srv.ServeTLS(l, tls.CertFilename, tls.KeyFilename)
		}
		srvFunc = wrapper
		srvMsg = "starting HTTPS server"
	}

	// Set up the Fx lifecycle controller
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info(srvMsg, zap.String("addr", srv.Addr))
			go func() {
				err := srvFunc(ln)
				if err != nil {
					log.Fatal("failed to start server", zap.Error(err))
				}
			}()
			return nil
		},

		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	// Return the http server
	return srv
}

// NewGinRouter returns a new instance of a *gin.Engine which is used to setup
// all of the request handling for the http server
func NewGinRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// NewController creates a new instance of the controller.Controller struct
func NewController(router *gin.Engine, logger *zap.Logger, cfg *config.Configuration,
	hrl *httpRoutes.RouteList, db *authndb.AuthDB) *controller.Controller {
	// Create the SubsystemParams object
	c := controller.Controller{
		Router:           router,
		Logger:           logger,
		Config:           cfg,
		HTTPRouteList:    hrl,
		SecureMiddleware: make([]gin.HandlerFunc, 0),
		AuthDB:           db,
	}
	return &c
}

// Register is a function that is called by the Fx framework to register all of the subsystems
func Register(params RegisterParams) {
	// Find all authentication middleware and setup
	for _, sub := range params.Subsystems {
		params.Controller.Logger.Info("checking for subsystem for authentication middleware", zap.String("subsystem", sub.Name()))
		if a, ok := sub.(interfaces.AuthenticationMiddleware); ok {
			params.Controller.Logger.Info("registering authentication middleware", zap.String("subsystem", sub.Name()))
			params.Controller.SecureMiddleware = append(params.Controller.SecureMiddleware, a.AuthenticationHandler)
		}
	}

	// Find all authorization middleware and setup
	for _, sub := range params.Subsystems {
		params.Controller.Logger.Info("checking for subsystem for authentication middleware", zap.String("subsystem", sub.Name()))
		if a, ok := sub.(interfaces.AuthorizationMiddleware); ok {
			params.Controller.Logger.Info("registering authentication middleware", zap.String("subsystem", sub.Name()))
			params.Controller.SecureMiddleware = append(params.Controller.SecureMiddleware, a.AuthorizationHandler)
		}
	}

	// Call register function on all subsystems
	for _, sub := range params.Subsystems {
		if err := sub.Register(); err != nil {
			params.Controller.Logger.Error("failed to register routes for subsystem", zap.String("subsystem", sub.Name()))
		}
	}
}

// AsSubsystem is a helper function that is used to annotate a function as a subsystem
func AsSubsystem(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(interfaces.Subsystem)),
		fx.ResultTags(`group:"subsystems"`),
	)
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
			NewHTTPServer,                           // Web Server
			NewGinRouter,                            // Web Request Router
			config.NewConfiguration,                 // Configuration
			zap.NewDevelopment,                      // Structured Logger
			basic.NewTLSInfo,                        // Tls Cert Handler
			httpRoutes.NewRouteList,                 // Http Routing List
			NewController,                           // Wrapper around common subsystem input components
			authndb.NewAuthDB,                       // Authentication database
			AsSubsystem(authn.NewSubsystem),         // Authentication Subsystem
			AsSubsystem(authz.NewSubsystem),         // Authorization Subsystem
			AsSubsystem(basic.NewEchoSubsystem),     // Demo Subsystem
			AsSubsystem(basic.NewHomePageSubsystem), // Homepage handler
			AsSubsystem(plugin.NewSubsystem),        // Plugin Subsystem
			AsSubsystem(network.NewSubsystem),       // Network Subsystem
			AsSubsystem(diag.NewSubsystem),          // Diagnostic Subsystem
			AsSubsystem(hardware.NewSubsystem),      // Hardware Subsystem
			AsSubsystem(system.NewSubsystem),        // System Subsystem
		),
		// Invoke any functions needed to initialize everything. The empty anonymous functions are
		// used to ensure that the providers that return that type are initialized.
		fx.Invoke(
			authorization.EnsureAuthzModel,  // Ensure we have at least a default authorization model
			authorization.EnsureAuthzPolicy, // Ensure we have at least a default authorization policy
			func(*http.Server) {},           // Set up the http server by requiring it for an anonymous function
			Register,                        // Register all the subsystems
		),
	).Run()
}
