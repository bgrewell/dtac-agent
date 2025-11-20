# Webhost Modules - Developer Guide

## Overview

Webhost modules are DTAC-managed executables that serve embedded static web UIs and expose internal proxy routes. They can run in two modes:

- **Standalone Mode**: For local development with CLI/config and stdio logging
- **Managed Mode**: Production deployment with gRPC control, token issuance, and event streaming

## Quick Start

### Building the Example Module

```bash
cd cmd/webhosts/dashboard
go build -o dashboard
```

### Running in Standalone Mode

```bash
./dashboard
```

The dashboard will start on `http://localhost:8080` by default.

Test the endpoints:
```bash
# Access the dashboard UI
curl http://localhost:8080/

# Check health status
curl http://localhost:8080/healthz
```

### Running in Managed Mode

Managed mode is activated by setting the `DTAC_WEBHOSTS` environment variable:

```bash
DTAC_WEBHOSTS=true ./dashboard
```

In managed mode, the module:
1. Starts a gRPC server
2. Prints a `CONNECT{{...}}` handshake for DTAC to discover it
3. Waits for DTAC to call `Init` with configuration
4. Runs with DTAC-provided settings (listen addr/port, proxies, permissions)

## Creating Your Own Webhost Module

### 1. Implement the WebhostModule Interface

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
    // Extract logger from context if available
    if logger, ok := ctx.Value("logger").(webhost.Logger); ok {
        m.logger = logger
    }
    // Perform initialization
    return nil
}

func (m *MyModule) Run(ctx context.Context) error {
    // Start your HTTP server here
    // Use m.config.Listen for address/port
    // Set up proxies using webhost.SetupProxies()
    // Block until ctx is cancelled
    return nil
}
```

### 2. Create Main Function

```go
func main() {
    module := &MyModule{}
    host := webhost.NewDefaultWebhostHost(module)
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle shutdown signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()
    
    if err := host.Run(ctx); err != nil {
        log.Fatalf("Host failed: %v", err)
    }
}
```

### 3. Embed Static Files (Optional)

```go
import "embed"

//go:embed static
var staticFiles embed.FS

// In your Run method:
staticFS, _ := fs.Sub(staticFiles, "static")
mux.Handle("/", http.FileServer(http.FS(staticFS)))
```

### 4. Configure Proxies

Proxies are configured in the DTAC config file (managed mode) or can be hardcoded (standalone):

```go
proxies := []webhost.ProxyConfig{
    {
        ID:        "metrics-proxy",
        FromPath:  "/api/metrics",
        TargetURL: "http://internal-metrics:9090",
        AuthType:  "jwt_from_dtac",
        Scopes:    []string{"metrics:read"},
        StripPrefix: true,
    },
}

// In your Run method:
webhost.SetupProxies(mux, proxies, tokenProvider, logger)
```

## Configuration

### DTAC Configuration (webhosts.yaml)

```yaml
webhosts:
  enabled: true
  entries:
    mymodule:
      enabled: true
      executable: /path/to/mymodule
      listen:
        addr: "127.0.0.1"
        port: 8080
        tls_profile: "default"  # optional
      permissions:
        roles: ["admin", "viewer"]
        allowed_scopes: ["metrics:read", "logs:read"]
        default_scopes: ["metrics:read"]
      proxies:
        - id: "metrics-proxy"
          from_path: "/api/metrics"
          target_url: "http://internal-metrics:9090"
          auth:
            type: "jwt_from_dtac"
            scopes: ["metrics:read"]
      config:  # optional module-specific config
        theme: "dark"
        refresh_interval: 30
```

### Standalone Configuration

For standalone mode, you can read configuration from:
- Command-line flags
- Environment variables
- Configuration files (JSON/YAML)

Example using flags:

```go
var (
    addr = flag.String("addr", "127.0.0.1", "Listen address")
    port = flag.Int("port", 8080, "Listen port")
)

func main() {
    flag.Parse()
    // Use *addr and *port in your InitConfig
}
```

## Security Model

### Proxy Authentication

Webhost modules can proxy requests to internal services without exposing credentials to the browser:

1. **Browser** → Module's HTTP server (no auth required for static assets)
2. **Module** → Internal service (with JWT from DTAC)

The `jwt_from_dtac` auth type:
- Module requests token via `TokenProvider.GetToken()`
- DTAC validates scopes against allowed_scopes
- DTAC issues signed JWT
- Module injects JWT in `Authorization: Bearer <token>` header

### Token Lifecycle

```go
// Get token provider from context (managed mode)
tokenProvider := ctx.Value("tokenProvider").(webhost.TokenProvider)

// Request token with scopes
token, expiresAt, err := tokenProvider.GetToken(ctx, []string{"metrics:read"}, "metrics-proxy")
if err != nil {
    // Handle error (token denied or provider unavailable)
}

// Token is cached and auto-refreshed by the SDK
```

### Permissions

Modules are assigned:
- **Roles**: High-level access groups (e.g., "admin", "viewer")
- **Allowed Scopes**: Specific permissions the module can request
- **Default Scopes**: Automatically granted if no scopes specified

DTAC enforces these via Casbin policies.

## Logging

### In Module Code

```go
logger.Info("Server started", map[string]string{
    "addr": addr,
    "port": fmt.Sprintf("%d", port),
})

logger.Error("Failed to connect", map[string]string{
    "error": err.Error(),
})
```

### Managed Mode
Logs are sent to DTAC via `EventStream` and forwarded to the centralized logging system (zap).

### Standalone Mode
Logs are written to stdout/stderr with timestamps.

## Best Practices

### 1. Graceful Shutdown

Always handle context cancellation:

```go
select {
case <-ctx.Done():
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return server.Shutdown(shutdownCtx)
case err := <-errChan:
    return err
}
```

### 2. Health Checks

Implement a `/healthz` endpoint:

```go
mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"status":"healthy"}`)
})
```

### 3. Error Handling

Log errors with context:

```go
if err != nil {
    logger.Error("Operation failed", map[string]string{
        "operation": "fetch_data",
        "error": err.Error(),
        "module": m.Name(),
    })
    return err
}
```

### 4. Timeouts

Set appropriate timeouts for HTTP server and clients:

```go
server := &http.Server{
    Addr:         addr,
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

## Testing

### Unit Tests

Test your module implementation:

```go
func TestModuleInit(t *testing.T) {
    module := &MyModule{}
    config := &webhost.InitConfig{
        ModuleName: "test",
        Listen: webhost.ListenConfig{
            Addr: "127.0.0.1",
            Port: 8080,
        },
    }
    
    ctx := context.Background()
    err := module.OnInit(ctx, config)
    if err != nil {
        t.Fatalf("OnInit failed: %v", err)
    }
}
```

### Integration Tests

Test the full module lifecycle:

```bash
# Start module in background
./mymodule &
MODULE_PID=$!

# Wait for startup
sleep 2

# Test endpoints
curl -f http://localhost:8080/healthz || exit 1
curl -f http://localhost:8080/ || exit 1

# Cleanup
kill $MODULE_PID
```

## Troubleshooting

### Module Won't Start

Check:
- Port is not already in use: `netstat -an | grep 8080`
- Permissions for binding to port (ports < 1024 require root)
- Embedded files are built into binary: `go build` (not `go run`)

### Proxies Not Working

Check:
- Target URL is accessible from module's network
- Token provider is available (managed mode only)
- Scopes are included in allowed_scopes
- Path mapping is correct (use StripPrefix if needed)

### Logs Not Appearing

In managed mode:
- Ensure EventStream is connected
- Check DTAC's log output for errors

In standalone mode:
- Logs go to stdout/stderr
- Redirect to file: `./mymodule > module.log 2>&1`

## Example: Full Module

See `cmd/webhosts/dashboard` for a complete working example that demonstrates:
- Embedded static files
- Health check endpoint
- Proxy configuration
- Both standalone and managed modes

## Next Steps

1. Build and test the example dashboard module
2. Create your own module based on the example
3. Configure DTAC to manage your module
4. Test token issuance and proxy authentication
5. Deploy to production

## References

- Epic document: `docs/webhost-module-epic.md`
- Proto definitions: `api/grpc/webhost.proto`
- SDK source: `pkg/webhost/`
- Example module: `cmd/webhosts/dashboard/`
