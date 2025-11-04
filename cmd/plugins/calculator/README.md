# Calculator Plugin

This plugin demonstrates plugin-to-plugin communication in the DTAC Agent.

## Features

- **Basic Calculations**: Perform add, subtract, multiply, and divide operations
- **Plugin-to-Plugin Communication**: Demonstrates calling another plugin (HelloPlugin)
- **Plugin Discovery**: List all currently loaded plugins

## Endpoints

### `/plugins/calculator/calculate`
Perform a basic calculation.

**Method**: POST

**Request Body**:
```json
{
  "operation": "add",    // add, subtract, multiply, divide
  "a": 5,
  "b": 3
}
```

**Response**:
```json
{
  "result": 8
}
```

### `/plugins/calculator/calculate_via_hello`
Perform a calculation while demonstrating plugin-to-plugin communication by calling the HelloPlugin first.

**Method**: POST

**Request Body**:
```json
{
  "operation": "multiply",
  "a": 4,
  "b": 7
}
```

**Response**:
```json
{
  "result": 28,
  "callee_plugin": "HelloPlugin"
}
```

### `/plugins/calculator/list_plugins`
List all currently loaded plugins.

**Method**: GET

**Response**:
```json
{
  "plugins": ["HelloPlugin", "CalculatorPlugin", "DockerPlugin"],
  "count": 3
}
```

## Building

```bash
# Build the plugin
go build -o bin/plugins/calculator.plugin cmd/plugins/calculator/main.go

# Build with debugging symbols
go build -gcflags="all=-N -l" -o bin/plugins/calculator.plugin cmd/plugins/calculator/main.go
```

## Configuration

Add to `/etc/dtac/config.yaml`:

```yaml
plugins:
  enabled: true
  plugin_dir: /opt/dtac/plugins
  entries:
    calculator:
      enabled: true
      hash: ""  # Optional SHA256 hash for verification
```

## Testing

```bash
# Get authentication token
TOKEN=$(sudo dtac token)

# Test basic calculation
curl -ks -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","a":10,"b":5}' \
  https://localhost:8180/plugins/calculator/calculate

# Test plugin-to-plugin communication
curl -ks -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"operation":"divide","a":20,"b":4}' \
  https://localhost:8180/plugins/calculator/calculate_via_hello

# List available plugins
curl -ks -H "Authorization: Bearer $TOKEN" \
  https://localhost:8180/plugins/calculator/list_plugins
```

## Plugin-to-Plugin Communication

This plugin demonstrates the new plugin-to-plugin communication feature. The `calculate_via_hello` endpoint:

1. Calls the HelloPlugin to get a greeting message
2. Logs the interaction
3. Performs the requested calculation
4. Returns the result along with information about the called plugin

This pattern can be used to build more complex plugin workflows where plugins collaborate to provide functionality.

## See Also

- [Plugin-to-Plugin Communication Documentation](../../docs/plugin-to-plugin-communication.md)
- [Hello Plugin](../hello/README.md) - A simple example plugin
- [Plugin Development Guide](../../README.md#plugin-development)
