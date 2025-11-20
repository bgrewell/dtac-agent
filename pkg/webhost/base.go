package webhost

import (
	"context"
	"time"
)

// InitConfig contains the configuration passed to a webhost module during initialization.
// This config is either provided by DTAC (managed mode) or constructed from CLI flags/config (standalone).
type InitConfig struct {
	ModuleName  string
	Listen      ListenConfig
	Proxies     []ProxyConfig
	Permissions Permissions
	Config      map[string]interface{} // Module-specific config
}

// ListenConfig specifies where the module's HTTP server should bind.
type ListenConfig struct {
	Addr        string // IP address (e.g., "127.0.0.1", "0.0.0.0")
	Port        int    // TCP port
	TLSProfile  string // Optional TLS profile name
	TLSCertPath string // Path to TLS certificate
	TLSKeyPath  string // Path to TLS key
}

// ProxyConfig defines a reverse-proxy route the module should expose.
type ProxyConfig struct {
	ID          string            // Unique identifier for this proxy
	FromPath    string            // Path prefix on module's HTTP server (e.g., "/api/metrics")
	TargetURL   string            // Upstream URL to proxy to (e.g., "http://internal-metrics:9090")
	AuthType    string            // Auth mechanism (e.g., "jwt_from_dtac", "none")
	Scopes      []string          // Permission scopes required
	StripPrefix bool              // Whether to strip FromPath before forwarding
	Headers     map[string]string // Additional headers to inject
}

// Permissions defines what roles and scopes the module is allowed to use.
type Permissions struct {
	Roles         []string // Roles assigned to this module
	AllowedScopes []string // Scopes this module can request tokens for
	DefaultScopes []string // Scopes granted by default
}

// WebhostModule is the interface that webhost modules must implement.
// The host calls these methods to manage the module lifecycle.
type WebhostModule interface {
	// Name returns the module's unique name.
	Name() string

	// OnInit is called after the host completes initialization.
	// The module receives its configuration and can perform setup.
	OnInit(ctx context.Context, config *InitConfig) error

	// Run starts the module's main logic (e.g., HTTP server).
	// It should block until ctx is cancelled or an error occurs.
	Run(ctx context.Context) error
}

// TokenProvider is an abstraction for obtaining access tokens.
// In managed mode, this requests tokens from DTAC via gRPC.
// In standalone mode, this returns mock tokens or errors.
type TokenProvider interface {
	// GetToken requests an access token with the specified scopes.
	// forProxy optionally identifies which proxy this token is for.
	// Returns the JWT token string, expiry time, and any error.
	GetToken(ctx context.Context, scopes []string, forProxy string) (token string, expiresAt time.Time, err error)
}

// LogLevel represents log severity levels.
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelFatal LogLevel = "FATAL"
)

// Logger is an abstraction for logging.
// In managed mode, this sends logs over EventStream.
// In standalone mode, this writes to stdout/stderr.
type Logger interface {
	// Log emits a log message with the specified level and fields.
	Log(level LogLevel, message string, fields map[string]string)

	// Debug logs a debug message.
	Debug(message string, fields map[string]string)

	// Info logs an info message.
	Info(message string, fields map[string]string)

	// Warn logs a warning message.
	Warn(message string, fields map[string]string)

	// Error logs an error message.
	Error(message string, fields map[string]string)

	// Fatal logs a fatal message.
	Fatal(message string, fields map[string]string)
}

// HostMode represents whether the module is running standalone or DTAC-managed.
type HostMode string

const (
	HostModeStandalone HostMode = "standalone"
	HostModeManaged    HostMode = "managed"
)
