package endpoint

// Func is the type for an endpoint function
type Func func(in *EndpointRequest) (out *EndpointResponse, err error)
