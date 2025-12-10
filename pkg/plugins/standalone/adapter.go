package standalone

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StandaloneRESTConfig contains minimal configuration for standalone REST adapter
type StandaloneRESTConfig struct {
	// Port is the HTTP port to listen on
	Port int

	// EnableTLS enables HTTPS if true
	EnableTLS bool

	// CertFile is the path to the TLS certificate file (required if EnableTLS is true)
	CertFile string

	// KeyFile is the path to the TLS key file (required if EnableTLS is true)
	KeyFile string

	// EnableCORS enables CORS middleware
	EnableCORS bool

	// AllowedOrigins lists allowed CORS origins (defaults to ["*"] if EnableCORS is true)
	AllowedOrigins []string

	// LogLevel sets the logging level ("debug", "info", "warn", "error")
	LogLevel string
}

// StandaloneRESTAdapter wraps the REST functionality for standalone plugin use
type StandaloneRESTAdapter struct {
	server    *http.Server
	router    *gin.Engine
	logger    *zap.Logger
	config    *StandaloneRESTConfig
	endpoints []*endpoint.Endpoint
}

// NewStandaloneRESTAdapter creates a new standalone REST adapter for a plugin
func NewStandaloneRESTAdapter(plugin plugins.Plugin, config *StandaloneRESTConfig) (*StandaloneRESTAdapter, error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	if config.Port == 0 {
		config.Port = 8080 // Default port
	}

	// Set up logger
	logger, err := createLogger(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Set up gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(ginZapLoggerMiddleware(logger))

	// Set up CORS if enabled
	if config.EnableCORS {
		corsConfig := cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}

		if len(config.AllowedOrigins) == 0 || (len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*") {
			corsConfig.AllowAllOrigins = true
		} else {
			corsConfig.AllowOrigins = config.AllowedOrigins
		}

		router.Use(cors.New(corsConfig))
		logger.Info("CORS middleware enabled", zap.Strings("allowed_origins", config.AllowedOrigins))
	}

	adapter := &StandaloneRESTAdapter{
		router: router,
		logger: logger,
		config: config,
	}

	// Register the plugin
	if err := adapter.registerPlugin(plugin); err != nil {
		return nil, fmt.Errorf("failed to register plugin: %w", err)
	}

	// Set up the HTTP server
	adapter.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}

	return adapter, nil
}

// Start starts the REST server
func (a *StandaloneRESTAdapter) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		return err
	}

	protocol := "HTTP"
	if a.config.EnableTLS {
		protocol = "HTTPS"
	}

	a.logger.Info(fmt.Sprintf("starting standalone REST %s server", protocol),
		zap.String("addr", a.server.Addr),
		zap.Int("endpoints", len(a.endpoints)))

	go func() {
		var err error
		if a.config.EnableTLS {
			err = a.server.ServeTLS(ln, a.config.CertFile, a.config.KeyFile)
		} else {
			err = a.server.Serve(ln)
		}
		if err != nil && err != http.ErrServerClosed {
			a.logger.Error("server error", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the REST server
func (a *StandaloneRESTAdapter) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

// registerPlugin registers a plugin's endpoints with the REST adapter
func (a *StandaloneRESTAdapter) registerPlugin(plugin plugins.Plugin) error {
	// Create a registration request
	req := &api.RegisterRequest{
		DefaultSecure: false, // Standalone plugins don't use auth by default
		Config:        "{}",  // Empty config
	}

	// Register the plugin
	resp := &api.RegisterResponse{}
	if err := plugin.Register(req, resp); err != nil {
		return fmt.Errorf("plugin registration failed: %w", err)
	}

	// Convert API endpoints to internal endpoints
	rootPath := plugin.RootPath()
	for _, apiEp := range resp.Endpoints {
		ep := convertAPIEndpointToEndpoint(apiEp, plugin, rootPath)
		a.endpoints = append(a.endpoints, ep)
		a.registerEndpoint(ep)
	}

	a.logger.Info("plugin registered",
		zap.String("plugin", plugin.Name()),
		zap.String("root_path", rootPath),
		zap.Int("endpoints", len(resp.Endpoints)))

	return nil
}

// registerEndpoint registers a single endpoint with the router
func (a *StandaloneRESTAdapter) registerEndpoint(ep *endpoint.Endpoint) {
	var method string
	switch ep.Action {
	case endpoint.ActionRead:
		method = http.MethodGet
	case endpoint.ActionWrite:
		method = http.MethodPut
	case endpoint.ActionCreate:
		method = http.MethodPost
	case endpoint.ActionDelete:
		method = http.MethodDelete
	default:
		a.logger.Error("invalid action", zap.String("action", ep.Action.String()))
		return
	}

	a.router.Handle(method, ep.Path, func(c *gin.Context) {
		in, err := createInputArgs(c)
		if err != nil {
			a.logger.Error("failed to create input args", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		out, err := ep.Function(in)
		if err != nil {
			a.logger.Error("failed to execute endpoint", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Set headers from response
		for headerKey, headerValues := range out.Headers {
			for _, headerValue := range headerValues {
				c.Header(headerKey, headerValue)
			}
		}

		// Write response
		c.Data(http.StatusOK, "application/json", out.Value)
	})

	a.logger.Debug("endpoint registered",
		zap.String("method", method),
		zap.String("path", ep.Path))
}

// convertAPIEndpointToEndpoint converts an API endpoint to an internal endpoint
func convertAPIEndpointToEndpoint(apiEp *api.PluginEndpoint, plugin plugins.Plugin, rootPath string) *endpoint.Endpoint {
	// Construct the full path
	path := fmt.Sprintf("/%s/%s", rootPath, apiEp.Path)

	// Convert action string to Action type
	var action endpoint.Action
	switch apiEp.Action {
	case "read", "READ":
		action = endpoint.ActionRead
	case "write", "WRITE":
		action = endpoint.ActionWrite
	case "create", "CREATE":
		action = endpoint.ActionCreate
	case "delete", "DELETE":
		action = endpoint.ActionDelete
	default:
		action = endpoint.ActionRead
	}

	// Create the endpoint with a wrapper function that calls the plugin
	ep := &endpoint.Endpoint{
		Path:        path,
		Action:      action,
		Description: apiEp.Description,
		Secure:      apiEp.Secure,
		AuthGroup:   apiEp.AuthGroup,
		Function: func(in *endpoint.Request) (*endpoint.Response, error) {
			// Call the plugin's method using just the endpoint path (not the full path)
			methodKey := fmt.Sprintf("%s:%s", action, apiEp.Path)
			return plugin.Call(methodKey, in)
		},
	}

	return ep
}

// createInputArgs creates an endpoint.Request from a gin.Context
func createInputArgs(ctx *gin.Context) (*endpoint.Request, error) {
	input := &endpoint.Request{
		Metadata:   make(map[string]string),
		Headers:    make(map[string][]string),
		Parameters: make(map[string][]string),
		Body:       nil,
	}

	// Populate headers
	for k, v := range ctx.Request.Header {
		input.Headers[k] = v
	}

	// Populate query parameters
	for k, v := range ctx.Request.URL.Query() {
		input.Parameters[k] = v
	}

	// Read request body
	body, err := ctx.GetRawData()
	if err != nil {
		return nil, err
	}
	input.Body = body

	return input, nil
}

// createLogger creates a zap logger based on the log level
func createLogger(level string) (*zap.Logger, error) {
	var zapLevel zap.AtomicLevel
	switch level {
	case "debug":
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info", "":
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config := zap.NewProductionConfig()
	config.Level = zapLevel
	config.DisableStacktrace = true
	config.Encoding = "console"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z0700")

	return config.Build()
}

// ginZapLoggerMiddleware is a custom middleware for Gin that uses Zap logger
func ginZapLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	}
}
