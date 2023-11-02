package endpoint

import (
	"context"
)

// InputArgs is the struct that captures all the inputs that are available to be sent to endpoints
// Context can hold things like headers
// Params is a map of key/value pairs where the value can be of any type
// ExpectedArgs helps the API understand what arguments are expected
type InputArgs struct {
	Context context.Context     `json:"-"`
	Headers map[string][]string `json:"headers,omitempty"`
	Params  map[string][]string `json:"params,omitempty"`
	Body    []byte              `json:"body,omitempty"`
}
