package rest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/basic"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"time"
)

// NewAdapter creates a new REST adapter
func NewAdapter(c *controller.Controller, tls *basic.TLSInfo) (adapter interfaces.APIAdapter, err error) {
	name := "rest"
	logger := c.Logger.With(zap.String("module", name))
	r := &Adapter{
		controller: c,
		router:     gin.Default(),
		logger:     logger,
		tls:        tls,
		name:       name,
		formatter:  helpers.NewJSONResponseFormatter(c.Config, logger),
	}
	return r, r.setup()
}

// Adapter is the REST API adapter
type Adapter struct {
	server     *http.Server
	router     *gin.Engine
	tls        *basic.TLSInfo
	controller *controller.Controller
	logger     *zap.Logger
	name       string
	srvMsg     string
	srvFunc    func(net.Listener) error
	formatter  helpers.ResponseFormatter
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

				// TODO: Handle verification of params, headers etc
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
	a.server = &http.Server{Addr: fmt.Sprintf(":%d", a.controller.Config.Listener.Port), Handler: a.router}

	// Set up the serve function
	a.srvFunc = a.server.Serve
	a.srvMsg = "starting HTTP server"
	if a.tls.Enabled {
		wrapper := func(l net.Listener) error {
			return a.server.ServeTLS(l, a.tls.CertFilename, a.tls.KeyFilename)
		}
		a.srvFunc = wrapper
		a.srvMsg = "starting HTTPS server"
	}

	return nil
}

func (a *Adapter) shim(method string, ep endpoint.Endpoint) {
	a.router.Handle(method, ep.Path, func(c *gin.Context) {
		in, err := a.createInputArgs(c)
		if err != nil {
			a.logger.Error("failed to create input args", zap.Error(err))
			a.formatter.WriteError(c, err)
			return
		}

		err = ep.ValidateArgs(in)
		if err != nil {
			a.logger.Error("failed to validate input args", zap.Error(err))
			a.formatter.WriteError(c, err)
			return
		}

		out, err := ep.Function(in)
		if err != nil {
			a.logger.Error("failed to execute endpoint", zap.Error(err))
			a.formatter.WriteError(c, err)
			return
		}

		et := out.Context.Value(types.ContextExecDuration).(time.Duration)
		a.formatter.WriteResponse(c, et, out.Value)
	})
}

func (a *Adapter) createInputArgs(ctx *gin.Context) (*endpoint.InputArgs, error) {
	input := &endpoint.InputArgs{
		Context: ctx.Request.Context(),
		Params:  make(map[string]interface{}),
		Body:    nil,
	}

	// Populate headers
	for k, v := range ctx.Request.Header {
		input.Context = context.WithValue(input.Context, k, v[0])
	}

	// Populate query parameters
	for k, v := range ctx.Request.URL.Query() {
		input.Params[k] = v[0]
	}

	// Read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	input.Body = bytes.NewReader(body)

	return input, nil
}
