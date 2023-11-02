package utility

import (
	"encoding/json"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
)

// ConvertToAPIInputArgs converts an endpoint InputArgs to an API InputArgs
func ConvertToAPIInputArgs(ea *endpoint.InputArgs) *api.InputArgs {
	headers := make(map[string]*api.StringList)
	for key, values := range ea.Headers {
		headers[key] = &api.StringList{Values: values}
	}

	params := make(map[string]*api.StringList)
	for key, values := range ea.Params {
		params[key] = &api.StringList{Values: values}
	}

	return &api.InputArgs{
		Headers: headers,
		Params:  params,
		Body:    ea.Body,
	}
}

// ConvertToEndpointInputArgs converts an API InputArgs to an endpoint InputArgs
func ConvertToEndpointInputArgs(ia *api.InputArgs) *endpoint.InputArgs {
	headers := make(map[string][]string)
	for key, values := range ia.Headers {
		headers[key] = values.Values
	}

	params := make(map[string][]string)
	for key, values := range ia.Params {
		params[key] = values.Values
	}

	return &endpoint.InputArgs{
		Headers: headers,
		Params:  params,
		Body:    ia.Body,
	}
}

// ConvertToAPIReturnVal converts an endpoint ReturnVal to an API ReturnVal
func ConvertToAPIReturnVal(ev *endpoint.ReturnVal) *api.ReturnVal {
	headers := make(map[string]*api.StringList)
	for key, values := range ev.Headers {
		headers[key] = &api.StringList{Values: values}
	}

	params := make(map[string]*api.StringList)
	for key, values := range ev.Params {
		params[key] = &api.StringList{Values: values}
	}

	valueJSON, err := json.Marshal(ev.Value)
	if err != nil {
		// Handle error. For now, we just return nil.
		return nil
	}

	return &api.ReturnVal{
		Headers: headers,
		Params:  params,
		Value:   string(valueJSON),
	}
}

// ConvertToEndpointReturnVal converts an API ReturnVal to an endpoint ReturnVal
func ConvertToEndpointReturnVal(rv *api.ReturnVal) *endpoint.ReturnVal {
	headers := make(map[string][]string)
	for key, value := range rv.Headers {
		values := make([]string, 0)
		values = append(values, value.Values...)
		headers[key] = values
	}

	params := make(map[string][]string)
	for key, value := range rv.Params {
		values := make([]string, 0)
		values = append(values, value.Values...)
		params[key] = values
	}

	var value interface{}
	if err := json.Unmarshal([]byte(rv.Value), &value); err != nil {
		// Handle error. For now, we just return nil.
		return nil
	}

	return &endpoint.ReturnVal{
		Headers: headers,
		Params:  params,
		Value:   value,
	}
}
