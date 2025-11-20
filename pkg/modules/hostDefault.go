package modules

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/modules/utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

// DefaultModuleHost is the default interface for the module host
type DefaultModuleHost struct {
	api.UnimplementedModuleServiceServer
	Module     Module
	Proto      string
	IP         string
	APIVersion string
	port       int
	encryptor  *utility.RPCEncryptor
	grpcServer *grpc.Server
}

// Register acts as a shim between the gRPC interface and the module interface. It handles conversion then calls the
// module's Register method.
func (mh *DefaultModuleHost) Register(ctx context.Context, request *api.ModuleRegisterRequest) (*api.ModuleRegisterResponse, error) {
	// Convert the config
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return nil, err
	}

	// Register the module
	response := &api.ModuleRegisterResponse{}
	err = mh.Module.Register(request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Call acts as a shim between the gRPC interface and the module interface. It handles conversion then calls the
// module's Call method.
func (mh *DefaultModuleHost) Call(ctx context.Context, request *api.EndpointRequestMessage) (*api.EndpointResponseMessage, error) {
	// Call the module
	in := utility.APIEndpointRequestToEndpointRequest(request.Request)
	ret, err := mh.Module.Call(request.Method, in)
	if err != nil {
		return nil, err
	}

	// Return result
	out := utility.EndpointResponseToAPIEndpointResponse(ret)
	return &api.EndpointResponseMessage{
		Id:       1,  // Unused
		Error:    "", // Unused - hold-over from previous implementation doesn't make sense here since we return an actual error object
		Response: out,
	}, nil
}

// LoggingStream acts as a shim between the gRPC interface and the module interface. It handles setting up the logging
// channel so the module can send structure logging messages back to the agent.
func (mh *DefaultModuleHost) LoggingStream(req *api.LoggingArgs, stream api.ModuleService_LoggingStreamServer) error {
	// Client calls in to set up the logging then the server uses the stream as a channel to send logging messages
	return mh.Module.LoggingStream(stream)
}

// RequestToken handles JWT token requests from modules
func (mh *DefaultModuleHost) RequestToken(ctx context.Context, request *api.TokenRequest) (*api.TokenResponse, error) {
	// For now, return an error - this will be implemented when token provisioning is integrated
	return nil, fmt.Errorf("token provisioning not yet implemented")
}

// RefreshToken handles JWT token refresh requests from modules
func (mh *DefaultModuleHost) RefreshToken(ctx context.Context, request *api.TokenRefreshRequest) (*api.TokenResponse, error) {
	// For now, return an error - this will be implemented when token provisioning is integrated
	return nil, fmt.Errorf("token refresh not yet implemented")
}

// Serve starts the module host
func (mh *DefaultModuleHost) Serve() error {
	// Verify that the ENV variable is set else exit with helpful message
	if os.Getenv("DTAC_MODULES") == "" {
		fmt.Println("============================ WARNING ============================")
		fmt.Println("This is a DTAC module and is not designed to be executed directly")
		fmt.Println("Please use the DTAC agent to load this module")
		fmt.Println("==================================================================")
		os.Exit(-1)
	}

	// Check for certificate and key files passed via ENV variables
	cert := os.Getenv("DTAC_TLS_CERT")
	key := os.Getenv("DTAC_TLS_KEY")

	// gRPC server setup
	var opts []grpc.ServerOption

	// Check if both certificate and key are provided
	if cert != "" && key != "" {
		// Create a certificate from the given cert and key
		tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
		if err != nil {
			return fmt.Errorf("failed to load server TLS certificate: %s", err)
		}

		// Create transport credentials for the gRPC server
		creds := credentials.NewServerTLSFromCert(&tlsCert)
		opts = append(opts, grpc.Creds(creds))
	}

	// Find a TCP port to use
	var err error
	mh.port, err = utility.GetUnusedTCPPort()
	if err != nil {
		return err
	}

	// gRPC setup
	mh.grpcServer = grpc.NewServer(opts...)
	api.RegisterModuleServiceServer(mh.grpcServer, mh)

	// Specify module rpc protocol
	rpcProto := "grpc"

	// Setup any options
	options := []string{
		fmt.Sprintf("enc=%s", url.QueryEscape(mh.encryptor.KeyString())), // Set up the option to enable encryption
		fmt.Sprintf("tls=%t", cert != "" && key != ""),                   // Set up the option to show if TLS is enabled
	}
	// Output connection information ( format: CONNECT{{NAME:ROOT_PATH:RPC_PROTO:TRANS_PROTO:IP:PORT:VER:OPTIONS}} )
	fmt.Printf("CONNECT{{%s:%s:%s:%s:%s:%d:%s:[%s]}}\n", mh.Module.Name(), mh.Module.RootPath(), rpcProto, mh.Proto, mh.IP, mh.port, mh.APIVersion, strings.Join(options, ","))

	// Listen for connections
	l, e := net.Listen(mh.Proto, fmt.Sprintf("%s:%d", mh.IP, mh.port))
	if e != nil {
		return e
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatalf("Failed to close listener: %v", err)
		}
	}(l)

	// Serve gRPC connections
	if err := mh.grpcServer.Serve(l); err != nil {
		return err
	}

	return nil
}

// GetPort returns the port the module host is listening on
func (mh *DefaultModuleHost) GetPort() int {
	return mh.port
}
