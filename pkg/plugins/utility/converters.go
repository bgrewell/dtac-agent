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
		Description:              ep.Description,
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
		Description:              ep.Description,
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

// EndpointRequestToAPIEndpointRequest converts an endpoint.Request to an api.Request
func EndpointRequestToAPIEndpointRequest(er *endpoint.Request) *api.EndpointRequest {
	headers := convertSliceToStringList(er.Headers)
	params := convertSliceToStringList(er.Parameters)

	return &api.EndpointRequest{
		Metadata:   er.Metadata,
		Headers:    headers,
		Parameters: params,
		Body:       er.Body,
	}
}

// APIEndpointRequestToEndpointRequest converts an api.Request to an endpoint.Request
func APIEndpointRequestToEndpointRequest(aer *api.EndpointRequest) *endpoint.Request {
	headers := convertStringListToSlice(aer.Headers)
	params := convertStringListToSlice(aer.Parameters)

	return &endpoint.Request{
		Metadata:   aer.Metadata,
		Headers:    headers,
		Parameters: params,
		Body:       aer.Body,
	}
}

// EndpointResponseToAPIEndpointResponse converts an endpoint.Response to an api.Response
func EndpointResponseToAPIEndpointResponse(er *endpoint.Response) *api.EndpointResponse {
	headers := convertSliceToStringList(er.Headers)
	params := convertSliceToStringList(er.Parameters)

	return &api.EndpointResponse{
		Metadata:   er.Metadata,
		Headers:    headers,
		Parameters: params,
		Value:      er.Value,
	}
}

// APIEndpointResponseToEndpointResponse converts an api.Response to an endpoint.Response
func APIEndpointResponseToEndpointResponse(aer *api.EndpointResponse) *endpoint.Response {
	headers := convertStringListToSlice(aer.Headers)
	params := convertStringListToSlice(aer.Parameters)

	return &endpoint.Response{
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
