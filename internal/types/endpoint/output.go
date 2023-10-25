package endpoint

import "context"

// ReturnVal is the struct that captures all the output that is returned to the endpoints
type ReturnVal struct {
	Context context.Context
	Value   interface{}
}
