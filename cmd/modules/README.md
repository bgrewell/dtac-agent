# DTAC Modules

DTAC modules are extensible components that can run either standalone or integrated with the DTAC agent.

## Running Modules

### Standalone Mode

Modules can run independently without the DTAC agent:

```bash
# Run a web module directly
./helloweb.module

# The web server will start and listen on the configured port
# Output: Web server listening on http://localhost:8090
```

In standalone mode:
- Web modules start their HTTP server directly
- Logging outputs to stdout with structured key=value format
- No RPC or gRPC communication is used
- Token provisioning is disabled

### DTAC Mode

When running under the DTAC agent, modules communicate via gRPC:

```bash
# Set the DTAC_MODULES environment variable
export DTAC_MODULES=1
./helloweb.module

# Output: CONNECT{{HelloWebModule:helloweb:grpc:tcp:127.0.0.1:PORT:...}}
```

In DTAC mode:
- Modules output a CONNECT message for the agent to parse
- gRPC server starts for bidirectional communication
- Logging uses RPC streaming to the agent
- Token provisioning is available (when implemented)

## Module Types

### Web Modules

Web modules serve static content and/or provide web-based interfaces:
- Embed static files using `go:embed`
- Implement the `WebModule` interface
- Start an HTTP server in standalone mode
- Example: `helloweb`

### API Modules

API modules provide REST/RPC endpoints:
- Require DTAC agent to expose endpoints
- Show informative message in standalone mode
- Example: `hello`

## Development

### Creating a New Module

1. Create a new directory under `cmd/modules/yourmodule`
2. Implement the `Module` interface from `pkg/modules`
3. For web modules, embed `WebModuleBase`
4. For API modules, embed `ModuleBase`
5. Add a `build.yaml` file for the build system

Example structure:
```
cmd/modules/yourmodule/
├── build.yaml
├── main.go
└── yourmodule/
    └── yourmodule.go
```

### Build Configuration

Create a `build.yaml` file:
```yaml
name: yourmodule
entry: main.go
platforms:
  - linux:amd64
  - darwin:amd64
  - windows:amd64
```

### Building

```bash
# Build all modules
mage modules

# Output: bin/modules/*.module
```

## Logging

Modules use structured logging with automatic format switching:

**Standalone mode** - stdout with key=value pairs:
```
2025/11/20 22:05:33 [INFO] web server started successfully {port=8090}
```

**DTAC mode** - RPC streaming to agent:
```go
m.Log(modules.LoggingLevelInfo, "message", map[string]string{
    "key": "value",
})
```

## Security

- Log fields are automatically escaped to prevent injection attacks
- Special characters (newlines, tabs, etc.) are escaped
- All modules pass CodeQL security scanning
