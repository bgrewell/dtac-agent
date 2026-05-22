package osqueryplugin

import (
	"encoding/json"
	"time"
)

// Config is the runtime configuration for the osquery plugin. Values are
// populated from the JSON config string supplied by the dtac-agent at
// Register() time. All fields are optional; zero values fall back to
// sensible defaults applied by withDefaults.
type Config struct {
	// BinaryPath overrides the bundled osqueryd binary location. If empty,
	// the plugin looks for `osquery-bin/osqueryd` next to its own executable.
	BinaryPath string `json:"binary_path,omitempty"`

	// ExtraFlags is appended verbatim to the osqueryd command line.
	ExtraFlags []string `json:"extra_flags,omitempty"`

	// SocketPath overrides the auto-generated unix socket location.
	SocketPath string `json:"socket_path,omitempty"`

	// QueryTimeout bounds the duration of a single SQL query.
	QueryTimeout time.Duration `json:"query_timeout,omitempty"`

	// StartupTimeout bounds the wait for the osqueryd extension socket to
	// become available after the daemon is spawned.
	StartupTimeout time.Duration `json:"startup_timeout,omitempty"`

	// MaxRows is the cap on rows returned by the generic /query endpoint.
	MaxRows int `json:"max_rows,omitempty"`
}

// parseConfig accepts the raw JSON config from RegisterRequest.Config and
// returns a Config with defaults filled in. An empty string is treated as
// "use all defaults".
func parseConfig(raw string) (Config, error) {
	var cfg Config
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
			return Config{}, err
		}
	}
	cfg.withDefaults()
	return cfg, nil
}

func (c *Config) withDefaults() {
	if c.QueryTimeout <= 0 {
		c.QueryTimeout = 10 * time.Second
	}
	if c.StartupTimeout <= 0 {
		c.StartupTimeout = 5 * time.Second
	}
	if c.MaxRows <= 0 {
		c.MaxRows = 10000
	}
}
