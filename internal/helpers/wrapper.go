package helpers

import (
	"context"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"time"
)

// HandleWrapper is a generic handler that is used to help add additional context and measurements to calls without
// requiring the duplication of this code into every handler.
func HandleWrapper(in *endpoint.InputArgs, f func() (interface{}, error), description string) (out *endpoint.ReturnVal, err error) {
	start := time.Now()
	value, err := f()
	if err != nil {
		return nil, err
	}
	response := types.AnnotatedStruct{
		Description: description,
		Value:       value,
	}
	duration := time.Since(start)
	ctx := context.WithValue(in.Context, types.ContextExecDuration, duration)
	out = &endpoint.ReturnVal{
		Context: ctx,
		Value:   response,
	}
	return out, nil
}
