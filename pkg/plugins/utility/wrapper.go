package utility

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
)

// PluginHandleWrapperWithHeaders is a generic handler that is used to help add additional context and measurements to calls without
// requiring the duplication of this code into every handler.
func PluginHandleWrapperWithHeaders(in *endpoint.InputArgs, f func() (headers map[string][]string, retval interface{}, err error), description string) (out *endpoint.ReturnVal, err error) {
	//start := time.Now()
	headers, value, err := f()
	if err != nil {
		return nil, err
	}
	response := value
	// TODO: Disabled for now, think more about if/how this is exposed later
	//response := types.AnnotatedStruct{
	//	Description: description,
	//	Value:       value,
	//}
	//duration := time.Since(start) //TODO: Can't pass contexts through the RPC layer, so this is disabled for now
	//ctx := context.WithValue(in.Context, types.ContextExecDuration, -1)
	out = &endpoint.ReturnVal{
		Context: nil,
		Headers: headers,
		Value:   response,
	}
	return out, nil
}

// PluginHandleWrapper is a generic handler that is used to help add additional context and measurements to calls without
// requiring the duplication of this code into every handler.
func PluginHandleWrapper(in *endpoint.InputArgs, f func() (retval interface{}, err error), description string) (out *endpoint.ReturnVal, err error) {
	// Define a new function that matches the signature of the function expected by HandleWrapperWithHeaders
	newFunc := func() (headers map[string][]string, retval interface{}, err error) {
		retval, err = f()
		// returning nil for headers map
		return nil, retval, err
	}

	return PluginHandleWrapperWithHeaders(in, newFunc, description)
}
