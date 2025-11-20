# Epic: Webhost Module Support

## Overview

Add support for "webhost modules": a new class of DTAC-managed executables that serve embedded static web UIs and expose internal proxy routes. They run in two modes:

1. **DTAC-managed**: gRPC control + token issuance + event/log stream
2. **Standalone**: local dev with CLI/config + stdio logging

## Goals

- Provide a gRPC control API (`WebhostService`) to initialize modules, issue/refresh tokens, and stream logs/events
- Let DTAC config fully control listen address/port, proxy definitions, and allowed permission scopes for each module
- Provide a Go SDK to minimize module boilerplate (mode detection, logging abstraction, token helper, proxy helper)
- Ensure modules can run standalone for local dev and DTAC-managed in production
- Enforce token requests via Casbin and issue signed JWTs when allowed

## Architecture

### gRPC API (`api/grpc/webhost.proto`)

**Service**: `WebhostService`

**RPCs**:
- `Init(WebhostInitRequest) → WebhostInitResponse`: DTAC provides listen address/port, proxy definitions, permission scopes, raw module config
- `RequestToken(TokenRequest) → TokenResponse`: module requests access tokens (JWTs) for proxies/upstream calls; responses include token string and expiry or error
- `EventStream(stream WebhostEvent) ↔ (stream WebhostEventAck)`: bidirectional stream for structured logs/events and acknowledgements

**Key Messages**:
- `ListenConfig`: address, port, TLS settings
- `ProxyConfig`: from_path, target_url, auth type + scopes
- `Permissions`: roles, allowed_scopes, default_scopes
- `TokenRequest`/`TokenResponse`: scope requests and JWT responses
- `WebhostEvent`: structured log/event messages

### Config Model

DTAC configuration exposes a new top-level section:

```yaml
webhosts:
  enabled: true
  entries:
    dashboard:
      enabled: true
      executable: /path/to/dashboard
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
      config:  # optional module-specific config (JSON)
        theme: "dark"
        refresh_interval: 30
```

DTAC owns listen addr/port and proxy target URLs in managed mode.

### SDK Components (`pkg/webhost`)

#### `base.go`
Core types and interfaces:
- `InitConfig`: configuration passed to module on init
- `ListenConfig`: listen address/port/TLS
- `ProxyConfig`: proxy route definitions
- `Permissions`: roles and scopes
- `WebhostModule`: interface modules must implement
- `TokenProvider`: abstraction for obtaining tokens

#### `host.go`
`DefaultWebhostHost` implementation:
- Detects standalone vs managed mode via environment variable
- In managed mode: starts gRPC server, prints CONNECT handshake, implements WebhostService server handlers
- In standalone mode: reads local flags/config, creates synthetic InitConfig, logs to stdio
- Wires module lifecycle (OnInit, Run)

#### `proxy.go`
Reverse-proxy helper:
- Creates HTTP handlers for proxy routes
- Injects JWTs via TokenProvider when required
- Handles path mapping and header preservation

### Subsystem (`internal/webhost`)

DTAC responsibilities:
- Discover webhost config entries at startup
- For each enabled entry:
  - Spawn executable with `DTAC_WEBHOSTS` env flag
  - Wait for `CONNECT{{...}}` handshake line to find gRPC endpoint
  - Dial module's WebhostService (respect TLS config)
  - Call `Init` with resolved config
  - Open `EventStream` for logs/events; forward to DTAC logging (zap) with module context
  - Provide token issuance backend that enforces Casbin policies
- Track process lifecycle, restart on failure, graceful shutdown

### Security Model

#### Proxy Behavior
- Browser hits module's HTTP server on configured listen addr/port (serves static assets)
- Module exposes internal proxy routes (e.g., `/api/metrics`) that forward to `target_url`
- Upstream requests decorated with JWTs obtained from DTAC (credentials never go to browser)
- Supported auth types: `jwt_from_dtac` (module requests token with configured scopes)

#### Token/Casbin Integration
When module calls `RequestToken(scopes)`:
1. DTAC identifies subject as `webhost:<module_name>`
2. Validate requested scopes are subset of `allowed_scopes` in config
3. Optionally run Casbin Enforce for fine-grained rules
4. If allowed, issue signed JWT (subject includes module id; scopes claim included; short-lived expiry)
5. If denied, return error via `TokenResponse`

Module should cache tokens and refresh near expiry (SDK implements caching and auto-refresh).

### Logging / Telemetry

- **DTAC-managed mode**: modules send structured logs/events over `EventStream`; DTAC forwards to zap with module metadata
- **Standalone mode**: SDK logs to stdout/stderr
- Module base routes logs to `EventStream` in managed mode; fallback to stdio in standalone

### Mode Comparison

#### Standalone Mode
- Module reads CLI flags or local config for listen addr/port and proxies
- No DTAC gRPC; logs to stdio; `TokenProvider` returns error or mock
- Use for local dev and debugging static assets and proxy behavior

#### Managed Mode
- Module starts gRPC control server, prints `CONNECT{{...}}` with gRPC details
- DTAC connects, calls `Init`, module accepts tokens and `EventStream` for logging
- Module binds HTTP server to listen addr/port provided by DTAC Init

## Deliverables

### Files Created
- `docs/webhost-module-epic.md` (this file)
- `api/grpc/webhost.proto` - gRPC service definitions
- `pkg/webhost/base.go` - core types and interfaces
- `pkg/webhost/host.go` - host implementation with mode detection
- `pkg/webhost/proxy.go` - reverse-proxy helper
- `cmd/webhosts/dashboard/main.go` - example module launcher
- `cmd/webhosts/dashboard/dashboard.go` - example module implementation
- `docs/webhost-modules.md` - developer guide

### Example Module: Dashboard

Located in `cmd/webhosts/dashboard`, demonstrates:
- Embedded static file serving
- Health endpoint (`/healthz`)
- Proxy route configuration
- Both standalone and managed operation

## Acceptance Criteria

- [ ] Epic doc file exists and matches design summary
- [ ] `webhost.proto` file present and properly formatted
- [ ] `pkg/webhost` contains `base.go`, `host.go`, `proxy.go` with starter interfaces and implementations
- [ ] `cmd/webhosts/dashboard` builds and runs standalone (serves assets and `/healthz`)
- [ ] DTAC webhost subsystem (skeleton) can spawn example module, detect CONNECT, call Init, and receive test event on EventStream
- [ ] `TokenRequest` enforcement returns tokens only for allowed scopes (unit tests)
- [ ] Developer docs show how to run example both standalone and under DTAC-managed flow

## Testing & Verification

### Unit Tests
- Config parsing
- Proxy creation
- Token request validation

### Build Tests
```bash
go build ./cmd/webhosts/dashboard
./dashboard --help  # standalone mode
```

### Integration Tests
1. Start local DTAC dev instance (or stub controller)
2. Spawn example module with `DTAC_WEBHOSTS=true`
3. Simulate CONNECT parsing and gRPC Init call
4. Verify EventStream logs arrive
5. Make test RequestToken and verify JWT structure & claims

### CI
- Run gofmt/golangci-lint
- Run unit tests
- Ensure repo compiles
- Proto generation step (if required)

## Future Enhancements

- Additional auth types (OAuth2, mTLS)
- Enhanced proxy error handling and observability
- Hot-reload of module config
- Module health checks and auto-restart
- WebSocket proxy support
- Rate limiting and circuit breakers

## Security Considerations

- Proxies cannot be used to exfiltrate secrets (DTAC owns target URLs and credentials)
- Casbin integration ensures fine-grained access control
- Short-lived JWTs with auto-refresh
- Module isolation (each module runs in separate process)
- TLS support for module-DTAC communication

## References

- Similar pattern to existing plugin system (`api/grpc/plugin.proto`)
- Follows DTAC's existing gRPC patterns and CONNECT handshake
- Integrates with existing Casbin authorization framework
