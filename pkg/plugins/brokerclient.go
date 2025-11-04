package plugins

import (
	"context"
	"crypto/tls"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

// BrokerClient provides a client interface for plugins to communicate with the agent's broker
type BrokerClient struct {
	conn    *grpc.ClientConn
	client  api.AgentServiceClient
	address string
}

// NewBrokerClient creates a new broker client that connects to the agent's broker service
func NewBrokerClient(brokerAddress string) (*BrokerClient, error) {
	if brokerAddress == "" {
		return nil, fmt.Errorf("broker address is empty")
	}

	// Setup security - check if TLS is enabled via environment variable
	var opts []grpc.DialOption
	if os.Getenv("DTAC_TLS_CERT") != "" && os.Getenv("DTAC_TLS_KEY") != "" {
		// TLS is enabled
		cert := os.Getenv("DTAC_TLS_CERT")
		key := os.Getenv("DTAC_TLS_KEY")
		tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
		if err != nil {
			return nil, fmt.Errorf("failed to load client TLS certificate: %w", err)
		}
		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Connect to the broker
	conn, err := grpc.Dial(brokerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker at %s: %w", brokerAddress, err)
	}

	client := api.NewAgentServiceClient(conn)

	return &BrokerClient{
		conn:    conn,
		client:  client,
		address: brokerAddress,
	}, nil
}

// CallPlugin calls another plugin through the broker
func (bc *BrokerClient) CallPlugin(pluginName string, method string, action endpoint.Action, request *endpoint.Request) (*endpoint.Response, error) {
	if bc.client == nil {
		return nil, fmt.Errorf("broker client is not connected")
	}

	// Convert request to API format
	apiRequest := utility.EndpointRequestToAPIEndpointRequest(request)

	// Make the call
	req := &api.PluginCallRequest{
		TargetPlugin: pluginName,
		Method:       method,
		Action:       action.String(),
		Request:      apiRequest,
	}

	resp, err := bc.client.CallPlugin(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("failed to call plugin %s: %w", pluginName, err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("plugin call error: %s", resp.Error)
	}

	// Convert response back to endpoint format
	response := utility.APIEndpointResponseToEndpointResponse(resp.Response)

	return response, nil
}

// ListPlugins returns a list of all loaded plugins
func (bc *BrokerClient) ListPlugins() ([]string, error) {
	if bc.client == nil {
		return nil, fmt.Errorf("broker client is not connected")
	}

	resp, err := bc.client.ListPlugins(context.Background(), &api.ListPluginsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list plugins: %w", err)
	}

	return resp.PluginNames, nil
}

// Close closes the broker client connection
func (bc *BrokerClient) Close() error {
	if bc.conn != nil {
		return bc.conn.Close()
	}
	return nil
}
