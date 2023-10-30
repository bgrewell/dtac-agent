package endpoint

type EndpointFunc func(in *InputArgs) (out *ReturnVal, err error)
