package endpoint

// WithMetadata sets the metadata option for the endpoint
func WithMetadata(metadata interface{}) Validators {
	return func(v *validationOptions) {
		v.metadata = metadata
	}
}

// WithHeaders sets the headers option for the endpoint
func WithHeaders(headers interface{}) Validators {
	return func(v *validationOptions) {
		v.headers = headers
	}
}

// WithParameters sets the parameters option for the endpoint
func WithParameters(parameters interface{}) Validators {
	return func(v *validationOptions) {
		v.parameters = parameters
	}
}

// WithBody sets the body option for the endpoint
func WithBody(body interface{}) Validators {
	return func(v *validationOptions) {
		v.body = body
	}
}

// WithOutput sets the output option for the endpoint
func WithOutput(output interface{}) Validators {
	return func(v *validationOptions) {
		v.output = output
	}
}

// Validators is a function that takes a validationOptions struct and sets the options for the endpoint
type Validators func(validator *validationOptions)

// validationOptions is a struct that holds the options for the endpoint
type validationOptions struct {
	metadata   interface{}
	headers    interface{}
	parameters interface{}
	body       interface{}
	output     interface{}
}
