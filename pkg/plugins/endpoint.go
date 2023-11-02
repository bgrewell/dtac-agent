package plugins

// PluginEndpoint is a struct that defines the endpoint for a plugin
type PluginEndpoint struct {
	Path           string      `json:"path" yaml:"path" toml:"path" mapstructure:"path"`
	Action         string      `json:"action" yaml:"action" toml:"action" mapstructure:"action"`
	UsesAuth       bool        `json:"uses_auth" yaml:"uses_auth" toml:"uses_auth" mapstructure:"uses_auth"`
	ExpectedArgs   interface{} `json:"expected_args,omitempty" yaml:"expected_args,omitempty" toml:"expected_args,omitempty" mapstructure:"expected_args,omitempty"`
	ExpectedBody   interface{} `json:"expected_body,omitempty" yaml:"expected_body,omitempty" toml:"expected_body,omitempty" mapstructure:"expected_body,omitempty"`
	ExpectedOutput interface{} `json:"expected_output,omitempty" yaml:"expected_output,omitempty" toml:"expected_output,omitempty" mapstructure:"expected_output,omitempty"`
}
