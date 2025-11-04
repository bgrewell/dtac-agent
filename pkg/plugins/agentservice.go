package plugins

import (
	"context"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
)

// AgentServiceImpl implements the AgentService gRPC interface
// This service is hosted by the agent and allows plugins to communicate with each other
type AgentServiceImpl struct {
	api.UnimplementedAgentServiceServer
	broker PluginBroker
}

// NewAgentService creates a new AgentService implementation
func NewAgentService(broker PluginBroker) *AgentServiceImpl {
	return &AgentServiceImpl{
		broker: broker,
	}
}

// CallPlugin handles plugin-to-plugin calls through the broker
func (s *AgentServiceImpl) CallPlugin(ctx context.Context, req *api.PluginCallRequest) (*api.PluginCallResponse, error) {
	// Parse the action
	action, err := endpoint.ParseAction(req.Action)
	if err != nil {
		return &api.PluginCallResponse{
			Error: fmt.Sprintf("invalid action '%s': %v", req.Action, err),
		}, nil
	}

	// Convert the request
	request := utility.APIEndpointRequestToEndpointRequest(req.Request)

	// Make the broker call
	response, err := s.broker.Call(req.TargetPlugin, req.Method, action, request)
	if err != nil {
		return &api.PluginCallResponse{
			Error: err.Error(),
		}, nil
	}

	// Convert the response
	apiResponse := utility.EndpointResponseToAPIEndpointResponse(response)

	return &api.PluginCallResponse{
		Response: apiResponse,
		Error:    "",
	}, nil
}

// ListPlugins returns a list of all loaded plugins
func (s *AgentServiceImpl) ListPlugins(ctx context.Context, req *api.ListPluginsRequest) (*api.ListPluginsResponse, error) {
	pluginNames := s.broker.ListPlugins()
	return &api.ListPluginsResponse{
		PluginNames: pluginNames,
	}, nil
}
