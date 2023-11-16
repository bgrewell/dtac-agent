package endpoint

import "context"

// ReturnVal is the struct that captures all the output that is returned to the endpoints
type ReturnVal struct {
	Context context.Context     `json:"-"`
	Headers map[string][]string `json:"headers,omitempty"`
	Params  map[string][]string `json:"params,omitempty"`
	Value   interface{}         `json:"value,omitempty"`
}
