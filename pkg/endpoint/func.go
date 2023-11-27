package endpoint

// Func is the type for an endpoint function
type Func func(in *Request) (out *Response, err error)
