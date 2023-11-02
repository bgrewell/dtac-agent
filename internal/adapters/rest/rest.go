package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/basic"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"time"
)

// NewAdapter creates a new REST adapter
func NewAdapter(c *controller.Controller, tls *map[string]basic.TLSInfo) (adapter interfaces.APIAdapter, err error) {
	// Check to see if the REST API is enabled. If not return an error that it is disabled
	if !c.Config.APIs.REST.Enabled {
		return nil, errors.New("rest api is not enabled")
	}

	name := "rest:api"
	logger := c.Logger.With(zap.String("module", name))
	r := &Adapter{
		controller: c,
		router:     gin.Default(),
		logger:     logger,
		tls:        tls,
		name:       name,
		formatter:  NewJSONResponseFormatter(c.Config, logger),
	}
	return r, r.setup()
}

// Adapter is the REST API adapter
type Adapter struct {
	server     *http.Server
	router     *gin.Engine
	tls        *map[string]basic.TLSInfo
	controller *controller.Controller
	logger     *zap.Logger
	name       string
	srvMsg     string
	srvFunc    func(net.Listener) error
	formatter  ResponseFormatter
}

// Name returns the name of the REST API adapter
func (a *Adapter) Name() string {
	return a.name
}

// Register registers the subsystems with the API adapter
func (a *Adapter) Register(subsystems []interfaces.Subsystem) (err error) {
	// Iterate over the subsystems and register each of the endpoints
	for _, subsystem := range subsystems {
		a.logger.Debug("registering subsystem", zap.String("subsystem", subsystem.Name()))
		if subsystem.Enabled() {
			// TODO: BUG - If this is called in multiple API adapters then our EndpointList in the controller which is
			//       universal is going to have duplicate endpoints. This should either be per API or should show
			//       API + endpoint
			a.controller.EndpointList.AddEndpoints(subsystem.Endpoints())
			for _, ep := range subsystem.Endpoints() {
				a.logger.Debug("registering endpoint", zap.String("path", ep.Path), zap.Any("action", ep.Action))
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
					return errors.New("invalid action")
				}

				a.shim(method, ep)
			}
		}
	}

	return nil
}

// Start starts the REST API adapter
func (a *Adapter) Start(ctx context.Context) (err error) {
	ln, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		return err
	}

	a.logger.Info(a.srvMsg, zap.String("addr", a.server.Addr))
	go func() {
		err := a.srvFunc(ln)
		if err != nil {
			a.logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the REST API adapter
func (a *Adapter) Stop(ctx context.Context) (err error) {
	return a.server.Shutdown(ctx)
}

func (a *Adapter) setup() (err error) {
	// Create a new http server
	a.server = &http.Server{Addr: fmt.Sprintf(":%d", a.controller.Config.APIs.REST.Port), Handler: a.router}

	// Set up the serve function
	a.srvFunc = a.server.Serve
	a.srvMsg = "starting HTTP server"
	if a.controller.Config.APIs.REST.TLS.Enabled {
		if cfg, ok := (*a.tls)[a.controller.Config.APIs.REST.TLS.Profile]; ok {
			wrapper := func(l net.Listener) error {
				return a.server.ServeTLS(l, cfg.CertFilename, cfg.KeyFilename)
			}
			a.srvFunc = wrapper
			a.srvMsg = "starting HTTPS server"
			return nil
		}

		return errors.New("tls profile not found")

	}

	return nil
}

func (a *Adapter) shim(method string, ep *endpoint.Endpoint) {
	a.router.Handle(method, ep.Path, func(c *gin.Context) {
		in, err := a.createInputArgs(c)
		if err != nil {
			a.logger.Error("failed to create input args", zap.Error(err))
			a.formatter.WriteError(c, err)
			return
		}

		// TODO: Look at moving validation somewhere central so API's don't have to do it
		err = ep.ValidateArgs(in)
		if err != nil {
			a.logger.Error("failed to validate input args", zap.Error(err))
			a.formatter.WriteError(c, err)
			return
		}

		// TODO: Look at moving validation somewhere central so API's don't have to do it
		// TODO: Look at deserialization of body and storage in a central manner to ease endpoint dup code
		err = ep.ValidateBody(in)
		if err != nil {
			a.logger.Error("failed to validate input body", zap.Error(err))
			a.formatter.WriteError(c, err)
			return
		}

		// If this is a secured endpoint check for the authorization header
		if ep.UsesAuth {
			auth := c.GetHeader("Authorization")
			if auth == "" {
				a.formatter.WriteUnauthorizedError(c, errors.New("authorization header is missing"))
				return
			}
			in.Context = context.WithValue(in.Context, types.ContextAuthHeader, auth)
		}

		// Add additional context
		in.Context = context.WithValue(in.Context, types.ContextResourceAction, ep.Action)
		in.Context = context.WithValue(in.Context, types.ContextResourcePath, ep.Path)

		out, err := ep.Function(in)
		if err != nil {
			a.logger.Error("failed to execute endpoint", zap.Error(err))
			a.formatter.WriteError(c, err)
			return
		}

		// Set headers from out.Headers into gin.Context
		for headerKey, headerValues := range out.Headers {
			for _, headerValue := range headerValues {
				c.Header(headerKey, headerValue)
			}
		}

		// If timing information is present then write it out
		if et, ok := out.Context.Value(types.ContextExecDuration).(time.Duration); ok {
			a.formatter.WriteResponse(c, et, out.Value)
		} else {
			a.formatter.WriteResponse(c, -1, out.Value)
		}

	})
}

func (a *Adapter) createInputArgs(ctx *gin.Context) (*endpoint.InputArgs, error) {
	input := &endpoint.InputArgs{
		Context: ctx.Request.Context(),
		Headers: make(map[string][]string),
		Params:  make(map[string][]string),
		Body:    nil,
	}

	// Populate headers
	for k, v := range ctx.Request.Header {
		input.Headers[k] = v
	}

	// Populate query parameters
	for k, v := range ctx.Request.URL.Query() {
		input.Params[k] = v
	}

	// Read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	input.Body = body

	return input, nil
}
