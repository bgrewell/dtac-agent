package endpoint

import (
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

// NewEndpoint creates a new instance of the Endpoint struct
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
	ExpectedMetadataSchema string `json:"-" mapstructure:"expected_metadata_schema,omitempty"`

	// ExpectedHeadersSchema defines the JSON Schema for the expected headers structure in the request.
	ExpectedHeadersSchema string `json:"-" mapstructure:"expected_headers_schema,omitempty"`

	// ExpectedParametersSchema defines the JSON Schema for the expected parameters structure in the request.
	ExpectedParametersSchema string `json:"-" mapstructure:"expected_parameters_schema,omitempty"`

	// ExpectedBodySchema defines the JSON Schema for the expected body structure in the request.
	ExpectedBodySchema string `json:"-" mapstructure:"expected_body_schema,omitempty"`

	// ExpectedOutputSchema defines the JSON Schema for the expected output structure in the response.
	ExpectedOutputSchema string `json:"-" mapstructure:"expected_output_schema,omitempty"`

	// ExpectedMetadataDescription is a output friendly representation of the expected metadata schema.
	ExpectedMetadataDescription json.RawMessage `json:"expected_metadata_schema,omitempty" yaml:"expected_metadata_schema,omitempty" toml:"expected_metadata_schema,omitempty"`
	// ExpectedHeadersDescription is a output friendly representation of the expected headers schema.
	ExpectedHeadersDescription json.RawMessage `json:"expected_headers_schema,omitempty" yaml:"expected_headers_schema,omitempty" toml:"expected_headers_schema,omitempty"`
	// ExpectedParametersDescription is a output friendly representation of the expected parameters schema.
	ExpectedParametersDescription json.RawMessage `json:"expected_parameters_schema,omitempty" yaml:"expected_parameters_schema,omitempty" toml:"expected_parameters_schema,omitempty"`
	// ExpectedBodyDescription is a output friendly representation of the expected body schema.
	ExpectedBodyDescription json.RawMessage `json:"expected_body_schema,omitempty" yaml:"expected_body_schema,omitempty" toml:"expected_body_schema,omitempty"`
	// ExpectedOutputDescription is a output friendly representation of the expected output schema.
	ExpectedOutputDescription json.RawMessage `json:"expected_output_schema,omitempty" yaml:"expected_output_schema,omitempty" toml:"expected_output_schema,omitempty"`
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
	if validators.body != nil {
		err := e.SetExpectedBodySchema(validators.body)
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
func (e *Endpoint) GenerateSchemaFromInterface(data interface{}, additionalProperties bool) (string, error) {
	schema := jsonschema.Reflect(data)
	if additionalProperties {
		for _, definition := range schema.Definitions {
			definition.AdditionalProperties = jsonschema.TrueSchema
		}
	}

	bytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetExpectedMetadataSchema sets the expected metadata schema from a given struct.
func (e *Endpoint) SetExpectedMetadataSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data, false)
	if err != nil {
		return err
	}
	e.ExpectedMetadataSchema = schema

	return json.Unmarshal([]byte(schema), &e.ExpectedMetadataDescription)
}

// SetExpectedHeadersSchema sets the expected headers schema from a given struct.
func (e *Endpoint) SetExpectedHeadersSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data, true)
	if err != nil {
		return err
	}
	e.ExpectedHeadersSchema = schema

	return json.Unmarshal([]byte(schema), &e.ExpectedHeadersDescription)
}

// SetExpectedParameterSchema sets the expected metadata schema from a given struct.
func (e *Endpoint) SetExpectedParameterSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data, true)
	if err != nil {
		return err
	}
	e.ExpectedParametersSchema = schema

	return json.Unmarshal([]byte(schema), &e.ExpectedParametersDescription)
}

// SetExpectedBodySchema sets the expected parameter schema from a given struct.
func (e *Endpoint) SetExpectedBodySchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data, false)
	if err != nil {
		return err
	}
	e.ExpectedBodySchema = schema

	return json.Unmarshal([]byte(schema), &e.ExpectedBodyDescription)
}

// SetExpectedOutputSchema sets the expected parameter schema from a given struct.
func (e *Endpoint) SetExpectedOutputSchema(data interface{}) error {
	schema, err := e.GenerateSchemaFromInterface(data, false)
	if err != nil {
		return err
	}
	e.ExpectedOutputSchema = schema

	return json.Unmarshal([]byte(schema), &e.ExpectedOutputDescription)
}

// ValidateRequest validates the request against the expected schemas
func (e *Endpoint) ValidateRequest(request *Request) error {
	if err := ValidateAgainstSchema(request.Metadata, e.ExpectedMetadataSchema); err != nil {
		return err
	}
	if err := ValidateAgainstSchema(request.Headers, e.ExpectedHeadersSchema); err != nil {
		return err
	}
	if err := ValidateAgainstSchema(request.Parameters, e.ExpectedParametersSchema); err != nil {
		return err
	}
	if err := ValidateAgainstSchema(request.Body, e.ExpectedBodySchema); err != nil {
		return err
	}

	return nil
}

// ValidateResponse validates the response against the expected output schema
func (e *Endpoint) ValidateResponse(response *Response) error {
	if err := ValidateAgainstSchema(response.Metadata, e.ExpectedMetadataSchema); err != nil {
		return err
	}
	if err := ValidateAgainstSchema(response.Headers, e.ExpectedHeadersSchema); err != nil {
		return err
	}
	if err := ValidateAgainstSchema(response.Parameters, e.ExpectedParametersSchema); err != nil {
		return err
	}
	if err := ValidateAgainstSchema(response.Value, e.ExpectedOutputSchema); err != nil {
		return err
	}

	return nil
}

// ValidateAgainstSchema validates the given data against the given schema
func ValidateAgainstSchema(data interface{}, schemaStr string) error {
	if schemaStr == "" {
		return nil // No schema provided
	}

	// Ensure data isn't a json string
	if input, ok := data.([]byte); ok {
		var tmp interface{}
		err := json.Unmarshal(input, &tmp)
		if err == nil {
			data = tmp
		} else {
			return err
		}
	}

	// Load schema and data into gojsonschema loaders
	schemaLoader := gojsonschema.NewStringLoader(schemaStr)
	documentLoader := gojsonschema.NewGoLoader(data)

	// Perform validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err // handle validation error
	}

	if !result.Valid() {
		// Collect and return errors if validation fails
		var validationErrors []string
		for _, err := range result.Errors() {
			validationErrors = append(validationErrors, err.String())
		}
		return fmt.Errorf("validation failed: %s", strings.Join(validationErrors, "; "))
	}

	return nil
}
