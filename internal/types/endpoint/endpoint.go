package endpoint

import (
	"fmt"
	"reflect"
	"strings"
)

// Endpoint is a struct that abstracts the API endpoints from the concrete API protocols that are available
type Endpoint struct {
	Path         string                                                `json:"path" yaml:"path" toml:"path" mapstructure:"path"`
	Action       Action                                                `json:"action" yaml:"action" toml:"action" mapstructure:"action"`
	Function     func(input *InputArgs) (output *ReturnVal, err error) `json:"-" yaml:"-" toml:"-" mapstructure:"-"`
	UsesAuth     bool                                                  `json:"uses_auth" yaml:"uses_auth" toml:"uses_auth" mapstructure:"uses_auth"`
	ExpectedArgs interface{}                                           `json:"expected_args,omitempty" yaml:"expected_args,omitempty" toml:"expected_args,omitempty" mapstructure:"expected_args,omitempty"`
	ExpectedBody interface{}                                           `json:"expected_body,omitempty" yaml:"expected_body,omitempty" toml:"expected_body,omitempty" mapstructure:"expected_body,omitempty"`
}

// ValidateArgs validates the arguments of the request against the expected arguments
func (e *Endpoint) ValidateArgs(in *InputArgs) error {
	if e.ExpectedArgs == nil {
		return nil
	}

	expectedArgs := reflect.TypeOf(e.ExpectedArgs)
	actualArgs := reflect.TypeOf(in.Params)

	if expectedArgs.Kind() != reflect.Struct {
		return fmt.Errorf("ExpectedArgs must be a struct")
	}

	if actualArgs.Kind() != reflect.Map {
		return fmt.Errorf("Params must be a map")
	}

	for i := 0; i < expectedArgs.NumField(); i++ {
		field := expectedArgs.Field(i)
		jsonTag := field.Tag.Get("json")

		if _, ok := in.Params[jsonTag]; !ok {
			if strings.Contains(jsonTag, "omitempty") {
				continue // Optional field is missing, but that's okay
			}

			return fmt.Errorf("Missing parameter: %s", jsonTag)
		}
	}

	return nil
}

// ValidateBody validates the body of the request against the expected body
func (e *Endpoint) ValidateBody(input *InputArgs) error {
	if e.ExpectedBody == nil {
		return nil
	}

	expectedBody := reflect.TypeOf(e.ExpectedBody)
	actualBody := reflect.TypeOf(input.Body)

	if expectedBody.Kind() == reflect.Ptr {
		expectedBody = expectedBody.Elem()
	}

	if actualBody.Kind() == reflect.Ptr {
		actualBody = actualBody.Elem()
	}

	if !actualBody.AssignableTo(expectedBody) {
		return fmt.Errorf("Invalid body type: expected %s, got %s", expectedBody.Name(), actualBody.Name())
	}

	return nil
}
