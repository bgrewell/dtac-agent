package utility

import (
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	pluginutil "github.com/bgrewell/dtac-agent/pkg/plugins/utility"
)

// ConvertPluginEndpointToEndpoint converts an api.PluginEndpoint to an endpoint.Endpoint
func ConvertPluginEndpointToEndpoint(ep *api.PluginEndpoint) *endpoint.Endpoint {
	return pluginutil.ConvertPluginEndpointToEndpoint(ep)
}

// ConvertEndpointToPluginEndpoint converts an endpoint.Endpoint to an api.PluginEndpoint
func ConvertEndpointToPluginEndpoint(ep *endpoint.Endpoint) *api.PluginEndpoint {
	return pluginutil.ConvertEndpointToPluginEndpoint(ep)
}

// APIEndpointRequestToEndpointRequest converts an api.Request to an endpoint.Request
func APIEndpointRequestToEndpointRequest(aer *api.EndpointRequest) *endpoint.Request {
	return pluginutil.APIEndpointRequestToEndpointRequest(aer)
}

// EndpointRequestToAPIEndpointRequest converts an endpoint.Request to an api.Request
func EndpointRequestToAPIEndpointRequest(er *endpoint.Request) *api.EndpointRequest {
	return pluginutil.EndpointRequestToAPIEndpointRequest(er)
}

// EndpointResponseToAPIEndpointResponse converts an endpoint.Response to an api.Response
func EndpointResponseToAPIEndpointResponse(er *endpoint.Response) *api.EndpointResponse {
	return pluginutil.EndpointResponseToAPIEndpointResponse(er)
}

// APIEndpointResponseToEndpointResponse converts an api.Response to an endpoint.Response
func APIEndpointResponseToEndpointResponse(aer *api.EndpointResponse) *endpoint.Response {
	return pluginutil.APIEndpointResponseToEndpointResponse(aer)
}
