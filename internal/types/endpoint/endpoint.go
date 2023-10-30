package endpoint

import (
	"fmt"
	"reflect"
	"strings"
)

// Endpoint is a struct that abstracts the API endpoints from the concrete API protocols that are available
type Endpoint struct {
	Path           string      `json:"path" yaml:"path" toml:"path" mapstructure:"path"`
	Action         Action      `json:"action" yaml:"action" toml:"action" mapstructure:"action"`
	Function       Func        `json:"-" yaml:"-" toml:"-" mapstructure:"-"`
	UsesAuth       bool        `json:"uses_auth" yaml:"uses_auth" toml:"uses_auth" mapstructure:"uses_auth"`
	ExpectedArgs   interface{} `json:"expected_args,omitempty" yaml:"expected_args,omitempty" toml:"expected_args,omitempty" mapstructure:"expected_args,omitempty"`
	ExpectedBody   interface{} `json:"expected_body,omitempty" yaml:"expected_body,omitempty" toml:"expected_body,omitempty" mapstructure:"expected_body,omitempty"`
	ExpectedOutput interface{} `json:"expected_output,omitempty" yaml:"expected_output,omitempty" toml:"expected_output,omitempty" mapstructure:"expected_output,omitempty"`
}

// ValidateArgs validates the arguments of the request against the expected arguments
func (e *Endpoint) ValidateArgs(in *InputArgs) error {
	if e.ExpectedArgs == nil {
		return nil
	}

	expectedArgs := reflect.TypeOf(e.ExpectedArgs)
	actualArgs := reflect.TypeOf(in.Params)

	if expectedArgs.Kind() != reflect.Struct {
		return fmt.Errorf("expectedArgs must be a struct")
	}

	if actualArgs.Kind() != reflect.Map {
		return fmt.Errorf("params must be a map")
	}

	for i := 0; i < expectedArgs.NumField(); i++ {
		field := expectedArgs.Field(i)
		jsonTag := field.Tag.Get("json")

		if _, ok := in.Params[jsonTag]; !ok {
			if strings.Contains(jsonTag, "omitempty") {
				continue // Optional field is missing, but that's okay
			}

			return fmt.Errorf("missing parameter: %s", jsonTag)
		}
	}

	return nil
}

// ValidateBody validates the body of the request against the expected body
func (e *Endpoint) ValidateBody(input *InputArgs) error {
	// Bypassed for now since at least current APIs perform their own validation
	return nil
}
