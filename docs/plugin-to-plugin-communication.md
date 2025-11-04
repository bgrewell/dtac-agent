# Plugin-to-Plugin Communication

## Overview

The DTAC Agent now supports plugin-to-plugin communication through a clean, agent-mediated broker pattern. This allows plugins to call methods on other loaded plugins while maintaining the agent's role as the central coordinator.

## Architecture

```
┌──────────┐         ┌─────────────────┐         ┌──────────┐
│ Plugin A │────────▶│  Agent (Broker) │────────▶│ Plugin B │
└──────────┘  gRPC   │                 │  gRPC   └──────────┘
                     │  - Routing      │
                     │  - Logging      │
                     │  - Security     │
                     └─────────────────┘
```

### Key Components

1. **AgentService (gRPC)**: A gRPC service hosted by the agent that provides:
   - `CallPlugin`: Route calls from one plugin to another
   - `ListPlugins`: List all currently loaded plugins

2. **PluginBroker**: Agent-side interface that:
   - Maintains registry of loaded plugins
   - Routes calls between plugins
   - Handles errors and plugin lifecycle

3. **BrokerClient**: Plugin-side helper that:
   - Connects to the agent's broker service
   - Provides convenient methods for calling other plugins
   - Manages connection lifecycle

## Usage

### Initializing Plugin-to-Plugin Communication

In your plugin's `Register` method, initialize the broker client:

```go
func (p *MyPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
    // Initialize broker client for plugin-to-plugin communication
    if request.BrokerAddress != "" {
        err := p.InitializeBrokerClient(request.BrokerAddress)
        if err != nil {
            p.Log(plugins.LevelWarning, "failed to initialize broker client", 
                map[string]string{"error": err.Error()})
        } else {
            p.Log(plugins.LevelInfo, "broker client initialized", 
                map[string]string{"address": request.BrokerAddress})
        }
    }
    
    // ... rest of registration
    return nil
}
```

### Calling Another Plugin

Use the `CallPlugin` method from `PluginBase`:

```go
func (p *MyPlugin) ProcessData(in *endpoint.Request) (*endpoint.Response, error) {
    // Call another plugin
    response, err := p.CallPlugin(
        "HelloPlugin",           // Target plugin name
        "hello",                 // Method to call
        endpoint.ActionRead,     // Action type
        &endpoint.Request{       // Request data
            Parameters: map[string][]string{
                "key": {"value"},
            },
        },
    )
    
    if err != nil {
        return nil, fmt.Errorf("failed to call HelloPlugin: %w", err)
    }
    
    // Process the response
    // ...
}
```

### Listing Available Plugins

```go
func (p *MyPlugin) GetPluginList(in *endpoint.Request) (*endpoint.Response, error) {
    pluginList, err := p.ListAvailablePlugins()
    if err != nil {
        return nil, fmt.Errorf("failed to list plugins: %w", err)
    }
    
    // pluginList is []string containing names of loaded plugins
    // ...
}
```

## Example: Calculator Plugin

The calculator plugin demonstrates plugin-to-plugin communication:

```go
// CalculateViaHello calls the hello plugin before performing a calculation
func (c *CalculatorPlugin) CalculateViaHello(in *endpoint.Request) (*endpoint.Response, error) {
    // Parse request
    var req CalculationRequest
    json.Unmarshal(in.Body, &req)
    
    // Call hello plugin
    helloResp, err := c.CallPlugin("HelloPlugin", "hello", endpoint.ActionRead, &endpoint.Request{})
    if err != nil {
        c.Log(plugins.LevelWarning, "hello plugin call failed", 
            map[string]string{"error": err.Error()})
    } else {
        c.Log(plugins.LevelInfo, "hello plugin responded", 
            map[string]string{"size": strconv.Itoa(len(helloResp.Value))})
    }
    
    // Perform calculation
    result := calculateResult(req)
    
    return &endpoint.Response{Value: result}, nil
}
```

## Configuration

No additional configuration is required. The broker is automatically started when the plugin subsystem initializes, and the broker address is passed to plugins during registration.

### Building the Example Plugin

```bash
# Build the calculator plugin
go build -o bin/plugins/calculator.plugin cmd/plugins/calculator/main.go

# Build with debugging symbols
go build -gcflags="all=-N -l" -o bin/plugins/calculator.plugin cmd/plugins/calculator/main.go
```

### Testing Plugin-to-Plugin Communication

1. Ensure both the hello and calculator plugins are loaded
2. Call the calculator plugin's `calculate_via_hello` endpoint:

```bash
# Get auth token
TOKEN=$(sudo dtac token)

# Call the endpoint
curl -ks -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","a":5,"b":3}' \
  https://localhost:8180/plugins/calculator/calculate_via_hello
```

Expected response:
```json
{
  "result": 8,
  "callee_plugin": "HelloPlugin"
}
```

## Security Considerations

- **Access Control**: All plugin-to-plugin calls go through the agent, allowing for centralized access control
- **Audit Trail**: All calls are logged by the agent for audit purposes
- **TLS Support**: Broker connections can use TLS (inherited from plugin TLS configuration)
- **Isolation**: Plugins remain isolated in separate processes

## Performance Considerations

- **Two-Hop Overhead**: Plugin-to-plugin calls involve two gRPC hops (caller → agent → target)
- **Acceptable Trade-off**: The overhead is minimal and acceptable given the benefits of centralized control
- **Optimization**: For high-frequency scenarios, consider caching data or using alternative designs

## Best Practices

1. **Error Handling**: Always handle errors from plugin calls gracefully
2. **Logging**: Log plugin-to-plugin calls for debugging and audit purposes
3. **Timeouts**: Consider adding timeouts for plugin calls
4. **Avoid Circular Dependencies**: Don't create circular call chains between plugins
5. **Optional Feature**: Make plugin-to-plugin communication optional - plugins should work without it

## Troubleshooting

### Broker Not Available

If you see "broker client is not initialized", ensure:
1. The agent's plugin subsystem is properly configured
2. The broker server started successfully (check logs)
3. You called `InitializeBrokerClient` during registration

### Plugin Not Found

If you get "plugin not found" errors:
1. Verify the target plugin is loaded (check agent logs)
2. Use correct plugin name (case-sensitive)
3. Call `ListAvailablePlugins()` to see what's loaded

### Connection Errors

If broker connection fails:
1. Check firewall settings
2. Verify TLS configuration matches between agent and plugins
3. Check agent logs for broker server startup errors

## API Reference

### PluginBase Methods

```go
// InitializeBrokerClient sets up the broker connection
func (p *PluginBase) InitializeBrokerClient(brokerAddress string) error

// CallPlugin calls another plugin through the broker
func (p *PluginBase) CallPlugin(
    pluginName string,
    method string,
    action endpoint.Action,
    request *endpoint.Request,
) (*endpoint.Response, error)

// ListAvailablePlugins returns names of all loaded plugins
func (p *PluginBase) ListAvailablePlugins() ([]string, error)
```

### gRPC Service Definition

```protobuf
service AgentService {
  rpc CallPlugin(PluginCallRequest) returns (PluginCallResponse);
  rpc ListPlugins(ListPluginsRequest) returns (ListPluginsResponse);
}
```

## Future Enhancements

Potential improvements to the plugin-to-plugin communication system:

1. **Rate Limiting**: Add rate limiting for plugin calls
2. **Call Tracing**: Implement distributed tracing for plugin call chains
3. **Metrics**: Add metrics for plugin-to-plugin communication
4. **Access Policies**: Fine-grained access control between specific plugins
5. **Async Calls**: Support for asynchronous plugin-to-plugin calls

## Conclusion

The plugin-to-plugin communication feature provides a clean, idiomatic way for plugins to collaborate while maintaining the agent's central coordination role. It follows Go and gRPC best practices and integrates seamlessly with the existing plugin architecture.
