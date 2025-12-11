package plugins

import (
	"os"
	"strconv"
	"strings"
)

// StandaloneConfig holds configuration for running a plugin in standalone mode
type StandaloneConfig struct {
	// Enabled indicates whether standalone mode is enabled
	Enabled bool
	// Protocol specifies http or https
	Protocol string
	// Port specifies the port to listen on
	Port int
	// TLSCertPath is the path to the TLS certificate file
	TLSCertPath string
	// TLSKeyPath is the path to the TLS key file
	TLSKeyPath string
	// Host specifies the host to bind to (default: 0.0.0.0)
	Host string
	// Config is the JSON configuration to pass to the plugin's Register method
	Config string
}

// StandaloneOption is a function type for configuring standalone mode
type StandaloneOption func(*StandaloneConfig)

// WithStandalone enables standalone mode
func WithStandalone() StandaloneOption {
	return func(c *StandaloneConfig) {
		c.Enabled = true
	}
}

// WithProtocol sets the protocol (http or https)
func WithProtocol(protocol string) StandaloneOption {
	return func(c *StandaloneConfig) {
		c.Protocol = protocol
	}
}

// WithPort sets the port
func WithPort(port int) StandaloneOption {
	return func(c *StandaloneConfig) {
		c.Port = port
	}
}

// WithTLS sets TLS certificate and key paths
func WithTLS(certPath, keyPath string) StandaloneOption {
	return func(c *StandaloneConfig) {
		c.TLSCertPath = certPath
		c.TLSKeyPath = keyPath
		if certPath != "" && keyPath != "" {
			c.Protocol = "https"
		}
	}
}

// WithHost sets the host to bind to
func WithHost(host string) StandaloneOption {
	return func(c *StandaloneConfig) {
		c.Host = host
	}
}

// WithConfig sets the plugin configuration JSON
func WithConfig(config string) StandaloneOption {
	return func(c *StandaloneConfig) {
		c.Config = config
	}
}

// NewStandaloneConfig creates a new StandaloneConfig with default values and applies options
func NewStandaloneConfig(opts ...StandaloneOption) *StandaloneConfig {
	// Start with defaults
	config := &StandaloneConfig{
		Enabled:  false,
		Protocol: "http",
		Port:     8080,
		Host:     "0.0.0.0",
		Config:   "{}",
	}

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	// Override with environment variables if present
	config.applyEnvVars()

	return config
}

// applyEnvVars reads environment variables and overrides config values
func (c *StandaloneConfig) applyEnvVars() {
	// Check if standalone mode is enabled via ENV
	if enabled := os.Getenv("DTAC_STANDALONE"); enabled != "" {
		c.Enabled = strings.ToLower(enabled) == "true" || enabled == "1"
	}

	// Protocol (http or https)
	if protocol := os.Getenv("DTAC_STANDALONE_PROTOCOL"); protocol != "" {
		c.Protocol = strings.ToLower(protocol)
	}

	// Port
	if portStr := os.Getenv("DTAC_STANDALONE_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			c.Port = port
		}
	}

	// TLS Certificate path
	if certPath := os.Getenv("DTAC_STANDALONE_TLS_CERT"); certPath != "" {
		c.TLSCertPath = certPath
	}

	// TLS Key path
	if keyPath := os.Getenv("DTAC_STANDALONE_TLS_KEY"); keyPath != "" {
		c.TLSKeyPath = keyPath
	}

	// Host
	if host := os.Getenv("DTAC_STANDALONE_HOST"); host != "" {
		c.Host = host
	}

	// If both TLS cert and key are set, upgrade to HTTPS
	if c.TLSCertPath != "" && c.TLSKeyPath != "" {
		c.Protocol = "https"
	}
}
