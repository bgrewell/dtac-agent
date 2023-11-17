package utility

import (
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
)

// ConvertPluginEndpointToEndpoint converts an api.PluginEndpoint to an endpoint.Endpoint
func ConvertPluginEndpointToEndpoint(ep *api.PluginEndpoint) *endpoint.Endpoint {
	action, _ := endpoint.ParseAction(ep.Action)
	eep := &endpoint.Endpoint{ // Create an endpoint endpoint (vs a plugin endpoint)
		Path:                     ep.Path,
		Action:                   action,
		Secure:                   ep.Secure,
		Function:                 nil,
		AuthGroup:                ep.AuthGroup,
		ExpectedMetadataSchema:   ep.ExpectedMetadataSchema,
		ExpectedHeadersSchema:    ep.ExpectedHeadersSchema,
		ExpectedParametersSchema: ep.ExpectedParametersSchema,
		ExpectedBodySchema:       ep.ExpectedBodySchema,
		ExpectedOutputSchema:     ep.ExpectedOutputSchema,
	}
	return eep
}

// ConvertEndpointToPluginEndpoint converts an endpoint.Endpoint to an api.PluginEndpoint
func ConvertEndpointToPluginEndpoint(ep *endpoint.Endpoint) *api.PluginEndpoint {
	aep := &api.PluginEndpoint{
		Path:                     ep.Path,
		Action:                   ep.Action.String(),
		Secure:                   ep.Secure,
		AuthGroup:                ep.AuthGroup,
		ExpectedMetadataSchema:   ep.ExpectedMetadataSchema,
		ExpectedHeadersSchema:    ep.ExpectedHeadersSchema,
		ExpectedParametersSchema: ep.ExpectedParametersSchema,
		ExpectedBodySchema:       ep.ExpectedBodySchema,
		ExpectedOutputSchema:     ep.ExpectedOutputSchema,
	}
	return aep
}

// EndpointRequestToAPIEndpointRequest converts an endpoint.EndpointRequest to an api.EndpointRequest
func EndpointRequestToAPIEndpointRequest(er *endpoint.EndpointRequest) *api.EndpointRequest {
	headers := convertSliceToStringList(er.Headers)
	params := convertSliceToStringList(er.Parameters)

	return &api.EndpointRequest{
		Metadata:   er.Metadata,
		Headers:    headers,
		Parameters: params,
		Body:       er.Body,
	}
}

// APIEndpointRequestToEndpointRequest converts an api.EndpointRequest to an endpoint.EndpointRequest
func APIEndpointRequestToEndpointRequest(aer *api.EndpointRequest) *endpoint.EndpointRequest {
	headers := convertStringListToSlice(aer.Headers)
	params := convertStringListToSlice(aer.Parameters)

	return &endpoint.EndpointRequest{
		Metadata:   aer.Metadata,
		Headers:    headers,
		Parameters: params,
		Body:       aer.Body,
	}
}

// EndpointResponseToAPIEndpointResponse converts an endpoint.EndpointResponse to an api.EndpointResponse
func EndpointResponseToAPIEndpointResponse(er *endpoint.EndpointResponse) *api.EndpointResponse {
	headers := convertSliceToStringList(er.Headers)
	params := convertSliceToStringList(er.Parameters)

	return &api.EndpointResponse{
		Metadata:   er.Metadata,
		Headers:    headers,
		Parameters: params,
		Value:      er.Value,
	}
}

// APIEndpointResponseToEndpointResponse converts an api.EndpointResponse to an endpoint.EndpointResponse
func APIEndpointResponseToEndpointResponse(aer *api.EndpointResponse) *endpoint.EndpointResponse {
	headers := convertStringListToSlice(aer.Headers)
	params := convertStringListToSlice(aer.Parameters)

	return &endpoint.EndpointResponse{
		Metadata:   aer.Metadata,
		Headers:    headers,
		Parameters: params,
		Value:      aer.Value,
	}
}

func convertSliceToStringList(headers map[string][]string) map[string]*api.StringList {
	output := make(map[string]*api.StringList)
	for key, values := range headers {
		output[key] = &api.StringList{Values: values}
	}
	return output
}

func convertStringListToSlice(headers map[string]*api.StringList) map[string][]string {
	output := make(map[string][]string)
	for key, values := range headers {
		output[key] = values.Values
	}
	return output
}
