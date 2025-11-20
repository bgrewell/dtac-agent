# Webhost SDK

Go SDK for building DTAC webhost modules.

## Overview

This package provides the building blocks for creating webhost modules that can run both standalone (for local development) and DTAC-managed (for production deployment).

## Components

### `base.go`

Core types and interfaces:
- `InitConfig`: Configuration passed to modules on initialization
- `ListenConfig`: HTTP server binding configuration
- `ProxyConfig`: Reverse-proxy route definitions
- `Permissions`: Roles and scopes
- `WebhostModule`: Interface that modules must implement
- `TokenProvider`: Abstraction for obtaining access tokens
- `Logger`: Logging interface

### `host.go`

`DefaultWebhostHost` implementation:
- Automatic mode detection (standalone vs managed)
- gRPC server setup for managed mode
- CONNECT handshake printing
- Logger and token provider wiring
- Module lifecycle management

### `proxy.go`

Reverse-proxy utilities:
- `ProxyHandler`: Creates HTTP handlers for proxy routes
- `SetupProxies`: Configures all proxy routes on an HTTP mux
- JWT injection via TokenProvider
- Path stripping and header manipulation
- Error handling and logging

## Usage

```go
package main

import (
    "context"
    "github.com/bgrewell/dtac-agent/pkg/webhost"
)

type MyModule struct {
    config *webhost.InitConfig
    logger webhost.Logger
}

func (m *MyModule) Name() string {
    return "mymodule"
}

func (m *MyModule) OnInit(ctx context.Context, config *webhost.InitConfig) error {
    m.config = config
    if logger, ok := ctx.Value("logger").(webhost.Logger); ok {
        m.logger = logger
    }
    return nil
}

func (m *MyModule) Run(ctx context.Context) error {
    // Start HTTP server using m.config.Listen
    // Set up proxies using webhost.SetupProxies()
    // Block until ctx is cancelled
    return nil
}

func main() {
    module := &MyModule{}
    host := webhost.NewDefaultWebhostHost(module)
    
    ctx := context.Background()
    if err := host.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Documentation

See `docs/webhost-modules.md` for the complete developer guide.

## Examples

See `cmd/webhosts/dashboard` for a working example module.
