# Standalone Hello Plugin

This is an example plugin that demonstrates how to run a DTAC plugin with its own REST API interface, independent of the DTAC agent framework.

## Features

This plugin can run in two modes:

1. **Standalone Mode**: Runs its own HTTP/HTTPS REST server with no middleware or authentication
2. **Traditional Mode**: Runs as a standard DTAC plugin via gRPC

## Building

```bash
# From the repository root
go build -o bin/standalone-hello ./cmd/plugins/standalone-hello
```

Or use the mage build system:

```bash
go run tools/mage/mage.go build
```

## Running

### Standalone Mode (with REST API)

Run the plugin with its own REST server:

```bash
./bin/standalone-hello -standalone
```

#### Options

```bash
./bin/standalone-hello -standalone \
  -port 8080 \
  -cors=true \
  -log-level info
```

Available flags:
- `-standalone` - Enable standalone mode
- `-port` - HTTP port (default: 8080)
- `-cors` - Enable CORS (default: true)
- `-tls` - Enable HTTPS (default: false)
- `-cert` - TLS certificate file (required if -tls is true)
- `-key` - TLS key file (required if -tls is true)
- `-log-level` - Logging level: debug, info, warn, error (default: info)

#### With TLS/HTTPS

```bash
./bin/standalone-hello -standalone \
  -port 8443 \
  -tls \
  -cert /path/to/cert.pem \
  -key /path/to/key.pem
```

### Traditional Mode (via DTAC agent)

Run as a standard DTAC plugin:

```bash
./bin/standalone-hello
```

This mode requires the plugin to be loaded by the DTAC agent.

## API Endpoints

When running in standalone mode, the following endpoints are available:

### GET /hello/message

Returns a hello message from the plugin.

**Example:**
```bash
curl http://localhost:8080/hello/message
```

**Response:**
```json
{
  "message": "Hello from standalone plugin! This plugin is running its own REST API without the DTAC agent.",
  "plugin": "StandaloneHelloPlugin",
  "mode": "standalone"
}
```

### POST /hello/echo

Echoes back the provided message.

**Example:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"message": "Test message"}' \
  http://localhost:8080/hello/echo
```

**Response:**
```json
{
  "message": "Test message",
  "plugin": "StandaloneHelloPlugin",
  "mode": "echo"
}
```

## Use Cases

This standalone mode is ideal for:

- **Development & Testing**: Quickly test plugin functionality without running the full DTAC agent
- **Microservices**: Deploy plugins as independent microservices
- **Simplified Deployments**: Run plugins in environments where the full DTAC agent is not needed
- **Prototyping**: Rapidly prototype new plugin functionality

## Security Considerations

⚠️ **Important**: Standalone mode does not include:
- Authentication
- Authorization
- Rate limiting
- Advanced input validation

This mode is designed for:
- Development and testing environments
- Trusted network environments
- Use cases where security is handled at the network/infrastructure level

For production deployments requiring security, use the plugin with the full DTAC agent framework.

## Architecture

```
┌─────────────────────────────────────┐
│   Standalone Hello Plugin           │
│  ┌────────────────────────────────┐ │
│  │  Plugin Logic                  │ │
│  │  - GetMessage()                │ │
│  │  - Echo()                      │ │
│  └────────────────────────────────┘ │
│            ↓                         │
│  ┌────────────────────────────────┐ │
│  │  Standalone REST Adapter       │ │
│  │  - Gin Router                  │ │
│  │  - HTTP Server                 │ │
│  │  - Optional CORS               │ │
│  │  - Optional TLS                │ │
│  └────────────────────────────────┘ │
└─────────────────────────────────────┘
            ↓
    HTTP/HTTPS Requests
```

## Comparison: Standalone vs Traditional

| Feature | Standalone Mode | Traditional Mode |
|---------|----------------|------------------|
| Transport | HTTP/REST | gRPC |
| Authentication | None | Yes (via DTAC agent) |
| Authorization | None | Yes (via DTAC agent) |
| Middleware | None | Yes (via DTAC agent) |
| Multi-plugin | No | Yes |
| Hot-reload | No | Yes (via DTAC agent) |
| Setup Complexity | Low | Medium |
| Best For | Dev/Test, Simple deployments | Production, Complex deployments |

## Creating Your Own Standalone Plugin

To create your own plugin that supports standalone mode:

1. Implement the standard `plugins.Plugin` interface
2. In your main.go, add standalone mode support:

```go
import (
    "github.com/bgrewell/dtac-agent/pkg/plugins/standalone"
)

func main() {
    plugin := myplugin.NewMyPlugin()
    
    config := &standalone.StandaloneRESTConfig{
        Port: 8080,
        EnableCORS: true,
        LogLevel: "info",
    }
    
    adapter, _ := standalone.NewStandaloneRESTAdapter(plugin, config)
    adapter.Start(context.Background())
    
    // Wait for shutdown signal...
}
```

See this example for a complete reference implementation.
