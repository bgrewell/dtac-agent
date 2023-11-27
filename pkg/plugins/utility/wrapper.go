package utility

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"time"
)

// PluginHandleWrapperWithHeaders is a generic handler that is used to help add additional context and measurements to calls without
// requiring the duplication of this code into every handler.
func PluginHandleWrapperWithHeaders(in *endpoint.Request, f func() (headers map[string][]string, retval []byte, err error), description string) (out *endpoint.Response, err error) {
	start := time.Now()
	headers, value, err := f()
	if err != nil {
		return nil, err
	}

	duration := time.Since(start) //TODO: Can't pass contexts through the RPC layer, so this is disabled for now
	//ctx := context.WithValue(in.Context, types.ContextExecDuration, -1)
	out = &endpoint.Response{
		Metadata: map[string]string{"duration": duration.String()},
		Headers:  headers,
		Value:    value,
	}
	return out, nil
}

// PluginHandleWrapper is a generic handler that is used to help add additional context and measurements to calls without
// requiring the duplication of this code into every handler.
func PluginHandleWrapper(in *endpoint.Request, f func() (retval []byte, err error), description string) (out *endpoint.Response, err error) {
	// Define a new function that matches the signature of the function expected by HandleWrapperWithHeaders
	newFunc := func() (headers map[string][]string, retval []byte, err error) {
		retval, err = f()
		// returning nil for headers map
		return nil, retval, err
	}

	return PluginHandleWrapperWithHeaders(in, newFunc, description)
}
