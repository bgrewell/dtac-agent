# Standalone REST Adapter for DTAC Plugins

This package provides a standalone REST adapter that allows DTAC plugins to be run independently with their own REST API interface, without requiring the full DTAC agent framework.

## Features

- ✅ Minimal dependencies - only requires the plugin and basic configuration
- ✅ No authentication/authorization middleware (as designed for standalone use)
- ✅ Optional TLS/HTTPS support
- ✅ Optional CORS support
- ✅ Built-in logging with configurable levels
- ✅ Self-contained - can be embedded in any plugin

## Usage

### Basic Example

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bgrewell/dtac-agent/cmd/plugins/hello/helloplugin"
	"github.com/bgrewell/dtac-agent/pkg/plugins/standalone"
)

func main() {
	// Create your plugin
	plugin := helloplugin.NewHelloPlugin()

	// Configure the standalone REST adapter
	config := &standalone.StandaloneRESTConfig{
		Port:       8080,
		EnableCORS: true,
		LogLevel:   "info",
	}

	// Create the adapter
	adapter, err := standalone.NewStandaloneRESTAdapter(plugin, config)
	if err != nil {
		log.Fatal(err)
	}

	// Start the server
	ctx := context.Background()
	if err := adapter.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	adapter.Stop(ctx)
}
```

### With TLS

```go
config := &standalone.StandaloneRESTConfig{
	Port:       8443,
	EnableTLS:  true,
	CertFile:   "/path/to/cert.pem",
	KeyFile:    "/path/to/key.pem",
	EnableCORS: true,
	LogLevel:   "debug",
}
```

### Custom CORS Origins

```go
config := &standalone.StandaloneRESTConfig{
	Port:           8080,
	EnableCORS:     true,
	AllowedOrigins: []string{"https://example.com", "https://app.example.com"},
	LogLevel:       "info",
}
```

## Configuration Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Port` | `int` | `8080` | HTTP/HTTPS port to listen on |
| `EnableTLS` | `bool` | `false` | Enable HTTPS |
| `CertFile` | `string` | `""` | Path to TLS certificate (required if EnableTLS=true) |
| `KeyFile` | `string` | `""` | Path to TLS private key (required if EnableTLS=true) |
| `EnableCORS` | `bool` | `false` | Enable CORS middleware |
| `AllowedOrigins` | `[]string` | `["*"]` | List of allowed CORS origins |
| `LogLevel` | `string` | `"info"` | Logging level: "debug", "info", "warn", "error" |

## API Endpoints

When a plugin is registered with the standalone adapter, its endpoints are exposed at:

```
http://localhost:PORT/{plugin_root_path}/{endpoint_path}
```

For example, if your plugin has:
- Root path: `hello`
- Endpoint: `message`
- Action: `READ` (GET)

The endpoint will be available at:
```
GET http://localhost:8080/hello/message
```

## Security Considerations

⚠️ **Important**: The standalone adapter is designed for development, testing, or trusted network environments. It does **not** include:
- Authentication
- Authorization
- Rate limiting
- Input validation (beyond what the plugin implements)

For production use cases requiring security, use the plugin with the full DTAC agent framework.

## Differences from Full DTAC Agent

| Feature | Standalone | Full Agent |
|---------|-----------|------------|
| Authentication | ❌ No | ✅ Yes |
| Authorization | ❌ No | ✅ Yes |
| Middleware Chain | ❌ No | ✅ Yes |
| TLS | ✅ Optional | ✅ Yes |
| CORS | ✅ Optional | ✅ Yes |
| Logging | ✅ Basic | ✅ Advanced |
| Plugin Hot-reload | ❌ No | ✅ Yes |
| Multi-plugin | ❌ No | ✅ Yes |

## Example Plugins

See `cmd/plugins/standalone-hello` for a complete working example.
