package utility

import (
	"context"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
)

// PlugFuncWrapperWithHeaders is a generic handler that is used to help add additional context and measurements to calls without
// requiring the duplication of this code into every handler.
func PlugFuncWrapperWithHeaders(in *endpoint.InputArgs, out *endpoint.ReturnVal, f func() (headers map[string][]string, retval interface{}, err error), description string) (err error) {
	//start := time.Now()
	headers, value, err := f()
	if err != nil {
		return err
	}
	response := value
	// TODO: Disabled for now, think more about if/how this is exposed later
	//response := types.AnnotatedStruct{
	//	Description: description,
	//	Value:       value,
	//}
	//duration := time.Since(start)
	//ctx := context.WithValue(in.Context, types.ContextExecDuration, duration)
	var ctx context.Context = nil // Have to clear context because it can't travel over an RPC - TODO: This should be redesigned
	out = &endpoint.ReturnVal{
		Context: ctx,
		Headers: headers,
		Value:   response,
	}
	return nil
}

// PlugFuncWrapper is a generic handler that is used to help add additional context and measurements to calls without
// requiring the duplication of this code into every handler.
func PlugFuncWrapper(in *endpoint.InputArgs, out *endpoint.ReturnVal, f func() (retval interface{}, err error), description string) (err error) {
	// Define a new function that matches the signature of the function expected by HandleWrapperWithHeaders
	newFunc := func() (headers map[string][]string, retval interface{}, err error) {
		retval, err = f()
		// returning nil for headers map
		return nil, retval, err
	}

	return PlugFuncWrapperWithHeaders(in, out, newFunc, description)
}
