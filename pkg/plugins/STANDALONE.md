# Plugin Standalone Mode

The DTAC plugin framework now supports running plugins in standalone mode with a REST API interface. This allows plugins to be executed directly without the DTAC agent, making them accessible via HTTP/HTTPS endpoints.

## Features

- **REST API**: Plugins expose their endpoints as REST endpoints
- **Flexible Configuration**: Configure via environment variables or code
- **TLS Support**: Optional HTTPS with certificate configuration
- **Protocol Agnostic**: Designed to support additional protocols in the future
- **Options Pattern**: Clean, composable configuration using functional options

## Configuration

### Environment Variables

Standalone mode can be configured using the following environment variables:

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DTAC_STANDALONE` | Enable standalone mode | `false` | `true` |
| `DTAC_STANDALONE_PROTOCOL` | Protocol (http or https) | `http` | `https` |
| `DTAC_STANDALONE_PORT` | Port to listen on | `8080` | `8443` |
| `DTAC_STANDALONE_HOST` | Host to bind to | `0.0.0.0` | `127.0.0.1` |
| `DTAC_STANDALONE_TLS_CERT` | Path to TLS certificate | - | `/path/to/cert.pem` |
| `DTAC_STANDALONE_TLS_KEY` | Path to TLS key | - | `/path/to/key.pem` |

### Programmatic Configuration

You can also configure standalone mode programmatically using the options pattern:

```go
import "github.com/bgrewell/dtac-agent/pkg/plugins"

func main() {
    p := myplugin.NewMyPlugin()

    // Enable standalone mode with default settings (uses ENV vars)
    h, err := plugins.NewPluginHost(p, plugins.WithStandalone())

    // Or with explicit configuration
    h, err := plugins.NewPluginHost(p,
        plugins.WithStandalone(),
        plugins.WithProtocol("http"),
        plugins.WithPort(8080),
        plugins.WithHost("0.0.0.0"),
    )

    // HTTPS with TLS
    h, err := plugins.NewPluginHost(p,
        plugins.WithStandalone(),
        plugins.WithPort(8443),
        plugins.WithTLS("/path/to/cert.pem", "/path/to/key.pem"),
    )

    err = h.Serve()
    if err != nil {
        log.Fatal(err)
    }
}
```

## Usage

### Running a Plugin in Standalone Mode

Using environment variables:

```bash
# HTTP on port 8080
DTAC_STANDALONE=true DTAC_STANDALONE_PORT=8080 ./my-plugin

# HTTPS with TLS
DTAC_STANDALONE=true \
DTAC_STANDALONE_PORT=8443 \
DTAC_STANDALONE_TLS_CERT=/path/to/cert.pem \
DTAC_STANDALONE_TLS_KEY=/path/to/key.pem \
./my-plugin
```

### REST Endpoint Mapping

Plugin endpoints are automatically mapped to REST endpoints based on their configuration:

- **Path**: `/{rootPath}/{endpoint}`
- **HTTP Method**: Maps to plugin action
  - `GET` → `read`
  - `POST` → `create`
  - `PUT` → `write`
  - `PATCH` → `write`
  - `DELETE` → `delete`

### Built-in Endpoints

All standalone plugins automatically include:

- `GET /health` - Health check endpoint

### Example

For a plugin with:
- Root path: `hello`
- Endpoint: `hello` (read action)

The REST endpoint would be:
```bash
GET http://localhost:8080/hello/hello
```

Response:
```json
{
  "message": "Hello, World!"
}
```

## Backward Compatibility

Standalone mode is completely optional and backward compatible:

- **Without standalone options**: Plugins work as before with gRPC
- **With standalone options**: Plugins run in REST mode
- **Environment variables**: Only applied when standalone mode is enabled

Existing plugins can continue to use:

```go
h, err := plugins.NewPluginHost(p)  // Uses gRPC mode
```

## Implementation Details

### Request/Response Conversion

The REST host automatically converts between HTTP requests/responses and the plugin's internal `endpoint.Request`/`endpoint.Response` format:

- HTTP headers → `Request.Headers`
- Query parameters → `Request.Parameters`
- Request body → `Request.Body`
- Response body → `Response.Value`
- Response headers → `Response.Headers`

### Plugin Registration

When a plugin starts in standalone mode:

1. The REST host is created with the configuration
2. The plugin's `Register()` method is called to initialize endpoints
3. HTTP routes are set up to handle requests
4. The HTTP server starts listening

### Error Handling

- Invalid HTTP methods return `405 Method Not Allowed`
- Unknown endpoints return `500 Internal Server Error` with error details
- Request parsing errors return `400 Bad Request`

## Future Enhancements

The standalone mode is designed to support additional protocols in the future:

- GraphQL
- gRPC (direct, not via agent)
- WebSocket
- Custom protocols

To add support for a new protocol, implement a new host type (similar to `RESTPluginHost`) and update `NewPluginHost()` to instantiate it based on configuration.
