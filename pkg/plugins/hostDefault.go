package plugins

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

// DefaultPluginHost is the default interface for the plugin host
type DefaultPluginHost struct {
	api.UnimplementedPluginServiceServer
	Plugin     Plugin
	Proto      string
	IP         string
	APIVersion string
	port       int
	encryptor  *utility.RPCEncryptor
	grpcServer *grpc.Server
}

// Register acts as a shim between the gRPC interface and the plugin interface. It handles conversion then calls the
// plugin's Register method.
func (ph *DefaultPluginHost) Register(ctx context.Context, request *api.RegisterArgs) (*api.RegisterReply, error) {
	// Convert the config
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return nil, err
	}

	// Register the plugin
	req := RegisterArgs{
		DefaultSecure: request.DefaultSecure,
		Config:        config,
	}
	rep := RegisterReply{}
	err = ph.Plugin.Register(req, &rep)
	if err != nil {
		return nil, err
	}

	eps := make([]*api.PluginEndpoint, 0)
	for _, ep := range rep.Endpoints {
		protoEp, err := convertEndpointToProto(ep)
		if err != nil {
			return nil, err
		}
		eps = append(eps, protoEp)
	}
	reply := &api.RegisterReply{
		Endpoints: eps,
	}
	return reply, nil
}

// Call acts as a shim between the gRPC interface and the plugin interface. It handles conversion then calls the
// plugin's Call method.
func (ph *DefaultPluginHost) Call(ctx context.Context, request *api.PluginRequest) (*api.PluginResponse, error) {
	// Convert the input args
	in := utility.ConvertToEndpointInputArgs(request.InputArgs)
	out, err := ph.Plugin.Call(request.Method, in)
	if err != nil {
		return nil, err
	}

	outAPI := utility.ConvertToAPIReturnVal(out)
	return &api.PluginResponse{
		Id:        1,  // Unused
		Error:     "", // Unused - hold-over from previous implementation doesn't make sense here since we return an actual error object
		ReturnVal: outAPI,
	}, nil
}

// Serve starts the plugin host
func (ph *DefaultPluginHost) Serve() error {
	// Hacky way to keep the net.rpc package from complaining about some method signatures
	logger := log.Default()
	logger.SetOutput(io.Discard)

	// Verify that the ENV variable is set else exit with helpful message
	if os.Getenv("DTAC_PLUGINS") == "" {
		fmt.Println("============================ WARNING ============================")
		fmt.Println("This is a DTAC plugin and is not designed to be executed directly")
		fmt.Println("Please use the DTAC agent to load this plugin")
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
	ph.port, err = utility.GetUnusedTCPPort()
	if err != nil {
		return err
	}

	// gRPC setup
	ph.grpcServer = grpc.NewServer(opts...)
	api.RegisterPluginServiceServer(ph.grpcServer, ph)

	// Specify plugin rpc protocol
	rpcProto := "grpc"

	// Setup any options
	options := []string{
		fmt.Sprintf("enc=%s", url.QueryEscape(ph.encryptor.KeyString())), // Set up the option to enable encryption
	}
	// Output connection information ( format: CONNECT{{NAME:ROOT_PATH:RPC_PROTO:TRANS_PROTO:IP:PORT:VER:OPTIONS}} )
	fmt.Printf("CONNECT{{%s:%s:%s:%s:%s:%d:%s:[%s]}}\n", ph.Plugin.Name(), ph.Plugin.RootPath(), rpcProto, ph.Proto, ph.IP, ph.port, ph.APIVersion, strings.Join(options, ","))

	// Listen for connections
	l, e := net.Listen(ph.Proto, fmt.Sprintf("%s:%d", ph.IP, ph.port))
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
	if err := ph.grpcServer.Serve(l); err != nil {
		return err
	}

	return nil
}

// GetPort returns the port the plugin host is listening on
func (ph *DefaultPluginHost) GetPort() int {
	return ph.port
}

func convertEndpointToProto(endpoint *PluginEndpoint) (*api.PluginEndpoint, error) {
	argsJSON, err := json.Marshal(endpoint.ExpectedArgs)
	if err != nil {
		return nil, err
	}
	bodyJSON, err := json.Marshal(endpoint.ExpectedBody)
	if err != nil {
		return nil, err
	}
	outputJSON, err := json.Marshal(endpoint.ExpectedOutput)
	if err != nil {
		return nil, err
	}
	return &api.PluginEndpoint{
		Path:           endpoint.Path,
		Action:         endpoint.Action,
		UsesAuth:       endpoint.UsesAuth,
		ExpectedArgs:   string(argsJSON),
		ExpectedBody:   string(bodyJSON),
		ExpectedOutput: string(outputJSON),
	}, nil
}
