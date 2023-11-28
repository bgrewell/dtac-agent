package grpc

import (
	"context"
	"errors"
	"fmt"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/basic"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
)

// NewAdapter creates a new gRPC adapter
func NewAdapter(c *controller.Controller, tls *map[string]basic.TLSInfo) (adapter interfaces.APIAdapter, err error) {
	// Check to see if the JSON-RPC API is enabled. If not return an error that it is disabled
	if !c.Config.APIs.GRPC.Enabled {
		return nil, errors.New("grpc api is not enabled")
	}

	// Setup logger
	name := "api/grpc"
	logger := c.Logger.With(zap.String("module", name))

	r := &Adapter{
		controller: c,
		logger:     logger,
		tls:        tls,
		name:       name,
		endpoints:  make(map[string]*endpoint.Endpoint),
	}
	return r, r.setup()
}

// Adapter is the gRPC API adapter
type Adapter struct {
	api.UnimplementedAdapterServiceServer
	server     *grpc.Server
	listener   net.Listener
	tls        *map[string]basic.TLSInfo
	controller *controller.Controller
	logger     *zap.Logger
	endpoints  map[string]*endpoint.Endpoint
	name       string
}

// Name returns the name of the gRPC API adapter
func (a *Adapter) Name() string {
	return a.name
}

// Register registers the subsystems with the API adapter
func (a *Adapter) Register(subsystems []interfaces.Subsystem) (err error) {
	// Iterate over the subsystems and register each of the endpoints
	for _, subsystem := range subsystems {
		a.logger.Debug("registering subsystem", zap.String("subsystem", subsystem.Name()))
		if subsystem.Enabled() {
			for _, ep := range subsystem.Endpoints() {
				a.logger.Debug("registering endpoint", zap.String("path", ep.Path), zap.Any("action", ep.Action))
				method := fmt.Sprintf("%s:%s", ep.Action, ep.Path)
				a.endpoints[method] = ep
			}
		}
	}

	return nil
}

// Start starts the gRPC API adapter
func (a *Adapter) Start(ctx context.Context) (err error) {
	a.logger.Info("starting gRPC API server", zap.String("addr", a.listener.Addr().String()))
	go func() {
		err := a.server.Serve(a.listener)
		if err != nil {
			a.logger.Fatal("failed to start gRPC API server", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the gRPC API adapter
func (a *Adapter) Stop(ctx context.Context) (err error) {
	a.server.GracefulStop()
	return nil
}

// List implements the List RPC
func (a *Adapter) List(ctx context.Context, in *api.ListRequest) (*api.ListResponse, error) {
	// Implement your logic here
	// For example, return a list of endpoints
	a.logger.Info("list request received", zap.Any("request", in))
	eps := make([]*api.PluginEndpoint, 0)
	for _, ep := range a.controller.EndpointList.Endpoints {
		eps = append(eps, utility.ConvertEndpointToPluginEndpoint(ep))
	}
	return &api.ListResponse{
		Endpoints: eps,
	}, nil
}

// Call implements the Call RPC
func (a *Adapter) Call(ctx context.Context, in *api.EndpointRequestMessage) (*api.EndpointResponseMessage, error) {
	// Implement your logic for the Call RPC
	a.logger.Info("call request received", zap.Any("request", in))

	// Ensure request and method have been passed
	if in == nil || in.GetMethod() == "" || in.GetRequest() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	// Get the input
	method := in.GetMethod()
	request := utility.APIEndpointRequestToEndpointRequest(in.GetRequest())

	// Ensure that we have the metadata map
	if request.Metadata == nil {
		request.Metadata = make(map[string]string)
	}

	if ep, ok := a.endpoints[method]; ok {
		request.Metadata[types.ContextResourceAction.String()] = ep.Action.String()
		request.Metadata[types.ContextResourcePath.String()] = ep.Path

		response, err := ep.Function(request)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		message := api.EndpointResponseMessage{Response: utility.EndpointResponseToAPIEndpointResponse(response)}
		return &message, nil
	}

	return nil, status.Error(codes.NotFound, "method not found")
}

func (a *Adapter) setup() (err error) {
	var opts []grpc.ServerOption

	// Check if TLS is enabled in the configuration
	if a.controller.Config.APIs.GRPC.TLS.Enabled {
		// Load server's certificate and private key
		if cfg, ok := (*a.tls)[a.controller.Config.APIs.GRPC.TLS.Profile]; ok {
			creds, err := credentials.NewServerTLSFromFile(cfg.CertFilename, cfg.KeyFilename)
			if err != nil {
				return fmt.Errorf("failed to load TLS keys: %v", err)
			}
			opts = append(opts, grpc.Creds(creds))
			a.logger.Debug("starting gRPC API server with TLS")
		} else {
			return errors.New("tls profile not found")
		}
	} else {
		a.logger.Debug("Starting gRPC API server without TLS")
	}

	// Create a gRPC server object
	a.server = grpc.NewServer(opts...)

	// Register server
	api.RegisterAdapterServiceServer(a.server, a)

	// If reflection is enabled, register the reflection service
	if a.controller.Config.APIs.GRPC.Reflection {
		a.logger.Debug("registering gRPC reflection service")
		// Register reflection service on gRPC server.
		reflection.Register(a.server)
	}

	// Create listener
	a.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", a.controller.Config.APIs.GRPC.Port))
	if err != nil {
		return err
	}

	return nil
}

// NOTE:
// Example testing from command line
//
// Login
// grpcurl -insecure -d '{"method": "create:auth/login", "request": {"metadata": {}, "headers": {}, "parameters": {}, "body": "eyJ1c2VybmFtZSI6ICJhZG1pbiIsICJwYXNzd29yZCI6ICJuZWVkX3RvX2dlbmVyYXRlX2FfcmFuZG9tX3Bhc3N3b3JkX29uX2luc3RhbGxfb3JfZmlyc3RfcnVuIn0K" }}' 127.0.0.1:8181 frontend.AdapterService.Call
//
// Call to secured diag/
// grpcurl -insecure -d '{"method": "read:diag/", "request": {"metadata": {"auth_header": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJmNzAyN2RiLTUyOGUtNGYwYS1iNzk2LTUwYjYxY2U0YTliYSIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTcwMTIwNzg3MCwidXNlcl9pZCI6MX0.XUFBR_PIrxmJ5QqJuPs7vDxpBxmfE4BhX93Jk4q6OAE"}, "headers": {}, "parameters": {}, "body": [] }}' 127.0.0.1:8181 frontend.AdapterService.Call | jq -r .response.value | base64 -d
