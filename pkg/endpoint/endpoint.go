package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
)

func NewEndpoint(path string, action Action, description string, function Func, secure bool, authGroup string, validators ...Validators) *Endpoint {
	ep := Endpoint{
		Path:        path,
		Action:      action,
		Description: description,
		Function:    function,
		Secure:      secure,
		AuthGroup:   authGroup,
	}

	// Handle any input/output validation options
	v := &validationOptions{}
	for _, validator := range validators {
		validator(v)
	}

	// Generate schemas
	ep.GenerateSchemas(v)

	return &ep
}

// Endpoint abstracts API endpoints from concrete API protocols, making it adaptable to different API styles like REST, gRPC, etc.
type Endpoint struct {
	// Path specifies the endpoint's unique path or identifier.
	Path string `json:"path" yaml:"path" toml:"path" mapstructure:"path"`

	// Action represents the type of operation this endpoint performs (e.g., GET, POST for REST).
	Action Action `json:"action" yaml:"action" toml:"action" mapstructure:"action"`

	// Function is the actual function to be executed when this endpoint is called.
	Function Func `json:"-" yaml:"-" toml:"-" mapstructure:"-"`

	// Description is a text based description of the endpoint that is shown in documentation and help output.
	Description string `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty" mapstructure:"description,omitempty"`

	// Secure indicates whether this endpoint requires authentication and authorization.
	Secure bool `json:"secure" yaml:"secure" toml:"secure" mapstructure:"secure"`

	// AuthGroup specifies the minimum authorization group required to access this endpoint.
	// Possible values might include 'Admin', 'Operator', 'User', 'Guest', etc.
	AuthGroup string `json:"auth_group,omitempty" yaml:"auth_group,omitempty" toml:"auth_group,omitempty" mapstructure:"auth_group,omitempty"`

	// ExpectedMetadataSchema defines the JSON Schema for the expected metadata structure in the request.
	ExpectedMetadataSchema string `json:"expected_metadata_schema,omitempty" yaml:"expected_metadata_schema,omitempty" toml:"expected_metadata_schema,omitempty" mapstructure:"expected_metadata_schema,omitempty"`

	// ExpectedHeadersSchema defines the JSON Schema for the expected headers structure in the request.
	ExpectedHeadersSchema string `json:"expected_headers_schema,omitempty" yaml:"expected_headers_schema,omitempty" toml:"expected_headers_schema,omitempty" mapstructure:"expected_headers_schema,omitempty"`

	// ExpectedParametersSchema defines the JSON Schema for the expected parameters structure in the request.
	ExpectedParametersSchema string `json:"expected_parameters_schema,omitempty" yaml:"expected_parameters_schema,omitempty" toml:"expected_parameters_schema,omitempty" mapstructure:"expected_parameters_schema,omitempty"`

	// ExpectedBodySchema defines the JSON Schema for the expected body structure in the request.
	ExpectedBodySchema string `json:"expected_body_schema,omitempty" yaml:"expected_body_schema,omitempty" toml:"expected_body_schema,omitempty" mapstructure:"expected_body_schema,omitempty"`

	// ExpectedOutputSchema defines the JSON Schema for the expected output structure in the response.
	ExpectedOutputSchema string `json:"expected_output_schema,omitempty" yaml:"expected_output_schema,omitempty" toml:"expected_output_schema,omitempty" mapstructure:"expected_output_schema,omitempty"`
}

// GenerateSchemas handles conversion for Go structs to JSON Schemas for input/output validation definition
func (e *Endpoint) GenerateSchemas(validators *validationOptions) {
	if validators.metadata != nil {
		err := e.SetExpectedMetadataSchema(validators.metadata)
		if err != nil {
			//TODO: Need to figure out how to handle logging here
			fmt.Printf("[!] ERROR: %v\n", err)
		}
	}
	if validators.headers != nil {
		err := e.SetExpectedHeadersSchema(validators.headers)
		if err != nil {
			//TODO: Need to figure out how to handle logging here
			fmt.Printf("[!] ERROR: %v\n", err)
		}
	}
	if validators.parameters != nil {
		err := e.SetExpectedParameterSchema(validators.parameters)
		if err != nil {
			//TODO: Need to figure out how to handle logging here
			fmt.Printf("[!] ERROR: %v\n", err)
		}
	}
	if validators.output != nil {
		err := e.SetExpectedOutputSchema(validators.output)
		if err != nil {
			//TODO: Need to figure out how to handle logging here
			fmt.Printf("[!] ERROR: %v\n", err)
		}
	}
}

// GenerateSchemaFromInterface generates a JSON Schema from a given interface{}.
// The interface should ideally be a struct for meaningful schema generation.
func (e *Endpoint) GenerateSchemaFromInterface(data interface{}) (string, error) {
	schema := jsonschema.Reflect(data)
	bytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetExpectedMetadataSchema sets the expected metadata schema from a given struct.
func (e *Endpoint) SetExpectedMetadataSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data)
	if err != nil {
		return err
	}
	e.ExpectedMetadataSchema = schema
	return nil
}

// SetExpectedHeadersSchema sets the expected headers schema from a given struct.
func (e *Endpoint) SetExpectedHeadersSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data)
	if err != nil {
		return err
	}
	e.ExpectedHeadersSchema = schema
	return nil
}

// SetExpectedParameterSchema sets the expected metadata schema from a given struct.
func (e *Endpoint) SetExpectedParameterSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data)
	if err != nil {
		return err
	}
	e.ExpectedMetadataSchema = schema
	return nil
}

// SetExpectedOutputSchema sets the expected parameter schema from a given struct.
func (e *Endpoint) SetExpectedOutputSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data)
	if err != nil {
		return err
	}
	e.ExpectedHeadersSchema = schema
	return nil
}

// ValidateArgs validates the arguments of the request against the expected arguments
func (e *Endpoint) ValidateArgs(request *EndpointRequest) error {
	//TODO: Return errors until implemented. Currently just return nil to allow development/testing to continue
	return nil
	//return errors.New("this method has not been implemented")
}

// ValidateBody validates the body of the request against the expected body
func (e *Endpoint) ValidateBody(request *EndpointRequest) error {
	//TODO: Return errors until implemented. Currently just return nil to allow development/testing to continue
	return nil
	//return errors.New("this method has not been implemented")
}
