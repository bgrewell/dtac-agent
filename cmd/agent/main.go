package main

import (
	"context"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/plugin"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/diag"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// NewHTTPServer creates the webserver that handles the requests sent to the DTAC agente
func NewHTTPServer(lc fx.Lifecycle, router *gin.Engine, log *zap.Logger, tls *helpers.TlsInfo) *http.Server {

	// Create a new http server
	srv := &http.Server{Addr: ":8180", Handler: router}

	// Setup the serve function
	srvFunc := srv.Serve
	srvMsg := "starting HTTP server"

	if tls.Enabled {
		wrapper := func(l net.Listener) error {
			return srv.ServeTLS(l, tls.CertFilename, tls.KeyFilename)
		}
		srvFunc = wrapper
		srvMsg = "starting HTTPS server"
	}

	// Setup the Fx lifecycle controller
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

func main() {
	fx.New(
		// Setup zap logger with the Fx framework
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		// Setup the providers
		fx.Provide(
			NewHTTPServer,           // Web Server
			NewGinRouter,            // Web Request Router
			helpers.NewTlsInfo,      // Tls Cert HAndler
			config.NewConfiguration, // Configuration
			plugin.NewSubsystem,     // Plugin Subsystem
			diag.NewSubsystem,       // Diagnostic Subsystem
			zap.NewExample,          // Structured Logger
		),
		// Invoke any functions needed to initialize everything. The empty anonymous functions are
		// used to ensure that the providers that return that type are initialized.
		fx.Invoke(
			func(*config.Configuration) {},
			func(*http.Server) {},
			func(*gin.Engine) {},
			func(*plugin.PluginSubsystem) {},
			func(*diag.DiagSubsystem) {},
			func(*helpers.TlsInfo) {},
		),
	).Run()
}
