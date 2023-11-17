package endpoint

func WithMetadata(metadata interface{}) Validators {
	return func(v *validationOptions) {
		v.metadata = metadata
	}
}

func WithHeaders(headers interface{}) Validators {
	return func(v *validationOptions) {
		v.headers = headers
	}
}

func WithParameters(parameters interface{}) Validators {
	return func(v *validationOptions) {
		v.parameters = parameters
	}
}

func WithBody(body interface{}) Validators {
	return func(v *validationOptions) {
		v.body = body
	}
}

func WithOutput(output interface{}) Validators {
	return func(v *validationOptions) {
		v.output = output
	}
}

type Validators func(validator *validationOptions)

type validationOptions struct {
	metadata   interface{}
	headers    interface{}
	parameters interface{}
	body       interface{}
	output     interface{}
}
