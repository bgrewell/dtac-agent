package endpoint

// Func is the type for an endpoint function
type Func func(in *InputArgs) (out *ReturnVal, err error)
