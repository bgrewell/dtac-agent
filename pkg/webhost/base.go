package webhost

import (
	"context"
	"time"
)

// InitConfig is provided to modules on Init.
type InitConfig struct {
	ModuleName  string
	Listen      ListenConfig
	Proxies     []ProxyConfig
	Permissions Permissions
	RawConfig   string // raw JSON module config if present
}

type ListenConfig struct {
	Addr       string
	Port       int
	TLS        bool
	TLSProfile string
}

type ProxyConfig struct {
	ID        string
	FromPath  string
	TargetURL string
	Auth      ProxyAuth
}

type ProxyAuth struct {
	Type   string   // e.g., "jwt_from_dtac"
	Scopes []string // requested scopes when getting tokens
}

type Permissions struct {
	Roles         []string
	AllowedScopes []string
	DefaultScopes []string
}

// WebhostModule is the interface modules implement.
type WebhostModule interface {
	Name() string
	// Called with the InitConfig from DTAC (or synthetic in standalone).
	OnInit(cfg *InitConfig) error
	// Run begins serving; block until context cancelled or server exits.
	Run(ctx context.Context) error
}

// TokenProvider is implemented by the host; modules should use it to request tokens.
// In DTAC-managed mode this will forward to DTAC via RPC. In standalone mode it may return an error or a mock token.
type TokenProvider interface {
	GetToken(ctx context.Context, scopes []string, forProxy string) (token string, expiresAt time.Time, err error)
}
