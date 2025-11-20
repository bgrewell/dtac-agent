# DTAC Module System

## Overview

The DTAC Module System provides a framework for extending DTAC with separately-managed processes that can host web frontends, expose REST APIs, run background services, or other extensible components. Modules complement the existing plugin system, which focuses on backend telemetry and control capabilities.

## Key Features

- **Process Isolation**: Modules run as separate processes managed by the DTAC agent
- **gRPC Communication**: Stdio-based RPC with encryption for secure inter-process communication
- **API Endpoint Registration**: Modules can register REST/gRPC endpoints with DTAC (like plugins)
- **Structured Logging**: Modules can send structured log messages back to the DTAC agent
- **Web Module Support**: Built-in support for hosting web frontends with embedded static assets
- **Configuration Management**: Agent pushes configuration to modules during registration
- **TLS Support**: Optional mutual TLS for secure communication

## Architecture

### Module Lifecycle

1. **Discovery**: Agent finds module executables (`.module`, `.module.exe`, or `.module.app`)
2. **Launch**: Agent spawns module process with environment variables
3. **Handshake**: Module outputs connection info via stdout in CONNECT{{}} format
4. **Connection**: Agent establishes gRPC connection to module
5. **Registration**: Agent calls Register() RPC, providing configuration
6. **Operation**: Module runs independently, logging back to agent
7. **Shutdown**: Agent can terminate module via cancel context

### Module Types

#### Basic Module
- Simple background process
- Can expose API endpoints through DTAC
- Structured logging
- Configuration from agent

#### Web Module
- HTTP server hosting static assets
- Embedded files using Go's embed.FS
- Can also expose API endpoints through DTAC
- Configurable port and routes
- Future: Proxy routes for backend APIs

## Creating a Module

### 1. Basic Module Structure

```go
package mymodule

import (
    "encoding/json"
    api "github.com/bgrewell/dtac-agent/api/grpc/go"
    "github.com/bgrewell/dtac-agent/pkg/modules"
    "reflect"
)

// Ensure interface compliance
var _ modules.Module = &MyModule{}

func NewMyModule() *MyModule {
    m := &MyModule{
        ModuleBase: modules.ModuleBase{},
    }
    m.SetRootPath("mymodule")
    return m
}

type MyModule struct {
    modules.ModuleBase
}

func (m MyModule) Name() string {
    t := reflect.TypeOf(m)
    return t.Name()
}

func (m *MyModule) Register(request *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error {
    *reply = api.ModuleRegisterResponse{
        ModuleType:   "basic",
        Capabilities: []string{"logging"},
    }

    // Parse config
    var config map[string]interface{}
    err := json.Unmarshal([]byte(request.Config), &config)
    if err != nil {
        return err
    }

    // Log registration
    m.Log(modules.LoggingLevelInfo, "module registered", map[string]string{
        "module_type": "basic",
    })

    return nil
}
```

### 2. Module Main Entry Point

```go
package main

import (
    "github.com/yourorg/yourmodule/mymodule"
    "log"
    "github.com/bgrewell/dtac-agent/pkg/modules"
)

func main() {
    m := mymodule.NewMyModule()

    h, err := modules.NewModuleHost(m)
    if err != nil {
        log.Fatal(err)
    }

    err = h.Serve()
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Web Module with Static Assets

```go
package mywebmodule

import (
    "embed"
    "encoding/json"
    api "github.com/bgrewell/dtac-agent/api/grpc/go"
    "github.com/bgrewell/dtac-agent/pkg/modules"
    "io/fs"
)

//go:embed static
var staticFiles embed.FS

func NewMyWebModule() *MyWebModule {
    m := &MyWebModule{
        WebModuleBase: modules.WebModuleBase{},
    }
    m.SetRootPath("myweb")
    m.SetConfig(modules.WebModuleConfig{
        Port:        8090,
        StaticPath:  "/",
        ProxyRoutes: []modules.ProxyRouteConfig{},
    })
    return m
}

type MyWebModule struct {
    modules.WebModuleBase
}

func (m *MyWebModule) GetStaticFiles() fs.FS {
    staticFS, _ := fs.Sub(staticFiles, "static")
    return staticFS
}

func (m *MyWebModule) Register(request *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error {
    *reply = api.ModuleRegisterResponse{
        ModuleType:   "web",
        Capabilities: []string{"static_files", "http_server"},
    }

    // Start web server
    err := m.Start()
    if err != nil {
        return err
    }

    m.Log(modules.LoggingLevelInfo, "web server started", map[string]string{
        "port": fmt.Sprintf("%d", m.GetPort()),
    })

    return nil
}
```

## Configuration

### Agent Configuration

Add to `config.yaml`:

```yaml
modules:
  enabled: true
  dir: /opt/dtac/modules/
  group: modules
  load_unconfigured: false
  tls:
    enabled: true
    profile: default
  entries:
    mymodule:
      enabled: true
      config:
        custom_option: "value"
      hash: ""  # Optional SHA256 hash for verification
      user: ""  # Future: run as specific user
```

### Module Configuration

Configuration is passed to the module's Register() method as JSON:

```go
var config map[string]interface{}
err := json.Unmarshal([]byte(request.Config), &config)
```

## Building Modules

### Build Command

```bash
go build -o mymodule.module ./cmd/mymodule
```

### Installation

1. Copy module binary to `/opt/dtac/modules/`
2. Update agent configuration
3. Restart agent or use dynamic loading (future feature)

## Logging

Modules can send structured logs back to the agent:

```go
m.Log(modules.LoggingLevelInfo, "operation completed", map[string]string{
    "duration": "100ms",
    "items": "42",
})
```

Log levels:
- `LoggingLevelDebug`
- `LoggingLevelInfo`
- `LoggingLevelWarning`
- `LoggingLevelError`
- `LoggingLevelFatal`

## Security

### File Permissions
- Module executables must be writable only by root or the process owner
- Agent verifies permissions before execution

### Hash Verification
- Optional SHA256 hash check in configuration
- Prevents execution of tampered modules

### TLS Communication
- Optional mutual TLS between agent and module
- Certificates managed by agent TLS profiles

### Encryption
- RPC communication encrypted with per-module symmetric keys
- Keys generated randomly at module startup

## Examples

See the following example modules:

1. **hello**: Basic module with API endpoint registration
   - Location: `cmd/modules/hello`
   - Features: Registration, logging, REST API endpoint
   - Demonstrates: How to expose endpoints through DTAC

2. **helloweb**: Web module with embedded HTML
   - Location: `cmd/modules/helloweb`
   - Features: Static file serving, embedded assets, HTTP server
   - Note: Can also register API endpoints if needed

## Future Enhancements

### JWT Token Provisioning (Planned)
Modules will be able to request JWT tokens from the agent:

```go
tokenResp, err := m.RPC.RequestToken(ctx, &api.TokenRequest{
    Scopes: []string{"read:api"},
    ExpiresIn: 3600,
})
```

### Proxy Routes (Planned)
Web modules will support authenticated proxy routes:

```go
ProxyRoutes: []modules.ProxyRouteConfig{
    {
        Path: "/api",
        Target: "http://backend:8080",
        StripPath: true,
        InjectToken: true,
    },
}
```

## Troubleshooting

### Module Not Starting

1. Check agent logs for errors
2. Verify file permissions (must not be world-writable)
3. Test module directly with `DTAC_MODULES=true ./mymodule.module`
4. Verify hash if configured

### Module Crashing

1. Check module logs in agent output
2. Verify configuration is valid JSON
3. Test with minimal configuration

### Connection Issues

1. Verify TLS configuration matches agent profile
2. Check firewall rules for localhost
3. Ensure port is not already in use

## API Reference

### Module Interface

```go
type Module interface {
    Name() string
    Register(args *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error
    Call(method string, args *endpoint.Request) (out *endpoint.Response, err error)
    RootPath() string
    LoggingStream(stream api.ModuleService_LoggingStreamServer) error
}
```

### ModuleBase

Provides default implementations and helper methods:
- `Call(method, args)` - Routes method calls to registered handlers
- `RegisterMethods(endpoints)` - Registers endpoint handlers
- `Log(level, message, fields)` - Send structured logs to agent
- `SetRootPath(path)` - Set module's root path
- `RootPath()` - Get module's root path

### WebModuleBase

Extends ModuleBase for web modules:
- `Start()` - Start HTTP server
- `Stop()` - Stop HTTP server
- `GetPort()` - Get listening port
- `SetConfig(config)` - Set web configuration
- `GetStaticFiles()` - Must be implemented by concrete type

## Contributing

When creating new module types or extending the module system:

1. Follow existing patterns from plugin system
2. Add comprehensive logging
3. Include example implementations
4. Update this documentation
5. Add tests where applicable

## License

Same as the DTAC agent project.

## Exposing API Endpoints

Modules can register API endpoints with DTAC, allowing them to expose REST/gRPC APIs that are managed by the agent's authentication, authorization, and routing infrastructure.

### Registering Endpoints

```go
func (m *MyModule) Register(request *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error {
    *reply = api.ModuleRegisterResponse{
        ModuleType:   "basic",
        Capabilities: []string{"api"},
        Endpoints:    make([]*api.PluginEndpoint, 0),
    }

    // Define endpoints
    authz := endpoint.AuthGroupAdmin.String()
    endpoints := []*endpoint.Endpoint{
        endpoint.NewEndpoint("users", endpoint.ActionRead, 
            "list users", m.ListUsers, request.DefaultSecure, authz),
        endpoint.NewEndpoint("users", endpoint.ActionCreate, 
            "create user", m.CreateUser, request.DefaultSecure, authz),
    }

    // Register methods (maps endpoints to handler functions)
    m.RegisterMethods(endpoints)

    // Convert and add to response
    for _, ep := range endpoints {
        aep := utility.ConvertEndpointToPluginEndpoint(ep)
        reply.Endpoints = append(reply.Endpoints, aep)
    }

    return nil
}

// Handler function
func (m *MyModule) ListUsers(in *endpoint.Request) (*endpoint.Response, error) {
    // Your logic here
    users := []string{"alice", "bob"}
    body, _ := json.Marshal(users)
    
    return &endpoint.Response{
        Headers: map[string][]string{"Content-Type": {"application/json"}},
        Value:   body,
    }, nil
}
```

### Endpoint Paths

Endpoints are automatically namespaced under the module's root path and the configured module group:
- Module root path: `hello`
- Module group: `modules`
- Endpoint: `users`
- Final path: `/modules/hello/users`

### Use Cases

This is particularly useful for web modules that need to:
- Manage user data via REST API
- Expose configuration endpoints
- Provide programmatic access to module functionality
- Integrate with the DTAC authentication/authorization system

Rather than building a separate REST API server in your module, you can register endpoints with DTAC and leverage its existing infrastructure for:
- JWT authentication
- Role-based authorization
- Request validation
- Error handling
- Logging and metrics

