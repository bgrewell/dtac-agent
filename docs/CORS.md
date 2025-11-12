# CORS Configuration Guide

## Overview
Cross-Origin Resource Sharing (CORS) support has been added to the DTAC Agent REST API adapter, allowing web frontends and browsers to interact with the API from different origins.

## Configuration

CORS can be configured through the `apis.rest.cors` section in your configuration file:

```yaml
apis:
  rest:
    enabled: true
    port: 8180
    cors:
      enabled: false  # Enable/disable CORS
      allowed_origins:
        - "*"  # Allow all origins (use with caution in production)
      allowed_methods:
        - GET
        - POST
        - PUT
        - DELETE
        - OPTIONS
      allowed_headers:
        - Origin
        - Content-Type
        - Accept
        - Authorization
      exposed_headers:
        - Content-Length
      allow_credentials: false
      max_age: 3600  # Preflight cache duration in seconds
```

## Configuration Options

### `enabled` (boolean)
- **Default**: `false`
- **Description**: Enable or disable CORS middleware
- **Recommendation**: Keep disabled unless you need cross-origin access

### `allowed_origins` (array of strings)
- **Default**: `["*"]`
- **Description**: List of origins that are allowed to access the API
- **Examples**:
  - `["*"]` - Allow all origins (not recommended for production)
  - `["https://dashboard.example.com"]` - Allow specific origin
  - `["https://dashboard.example.com", "https://app.example.com"]` - Multiple origins
- **Recommendation**: Use specific origins in production

### `allowed_methods` (array of strings)
- **Default**: `["GET", "POST", "PUT", "DELETE", "OPTIONS"]`
- **Description**: HTTP methods that are allowed in CORS requests
- **Common values**: GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD

### `allowed_headers` (array of strings)
- **Default**: `["Origin", "Content-Type", "Accept", "Authorization"]`
- **Description**: HTTP headers that clients can use in actual requests
- **Common additions**: X-Requested-With, X-Custom-Header

### `exposed_headers` (array of strings)
- **Default**: `["Content-Length"]`
- **Description**: Headers that browsers are allowed to access
- **Use case**: Custom response headers that your frontend needs to read

### `allow_credentials` (boolean)
- **Default**: `false`
- **Description**: Allow cookies and authentication headers in CORS requests
- **Note**: Cannot be used with wildcard (`*`) origins
- **Recommendation**: Set to `true` only if you need authenticated CORS requests

### `max_age` (integer)
- **Default**: `3600` (1 hour)
- **Description**: How long (in seconds) browsers can cache preflight responses
- **Range**: Typically 0-86400 (24 hours)
- **Recommendation**: Higher values reduce preflight requests but may delay configuration updates

## Security Considerations

### Production Best Practices

1. **Use Specific Origins**: Never use `["*"]` in production
   ```yaml
   allowed_origins:
     - "https://your-frontend.example.com"
   ```

2. **Limit Methods**: Only allow methods your API actually supports
   ```yaml
   allowed_methods:
     - GET
     - POST
   ```

3. **Restrict Headers**: Only allow headers your API needs
   ```yaml
   allowed_headers:
     - Content-Type
     - Authorization
   ```

4. **Credentials**: Only enable if absolutely necessary
   ```yaml
   allow_credentials: true
   allowed_origins:
     - "https://trusted-frontend.example.com"  # Must be specific, not "*"
   ```

## Example Configurations

### Development Environment
```yaml
apis:
  rest:
    cors:
      enabled: true
      allowed_origins:
        - "*"
      allow_credentials: false
```

### Production Environment
```yaml
apis:
  rest:
    cors:
      enabled: true
      allowed_origins:
        - "https://dashboard.example.com"
        - "https://app.example.com"
      allowed_methods:
        - GET
        - POST
        - PUT
        - DELETE
      allowed_headers:
        - Content-Type
        - Authorization
      allow_credentials: true
      max_age: 7200
```

### API Gateway / Proxy Setup
```yaml
apis:
  rest:
    cors:
      enabled: true
      allowed_origins:
        - "https://api-gateway.example.com"
      allowed_methods:
        - GET
        - POST
      allow_credentials: false
      max_age: 3600
```

## Testing CORS

### Using curl
```bash
# Test preflight request
curl -X OPTIONS http://localhost:8180/api/endpoint \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v

# Test actual request
curl -X GET http://localhost:8180/api/endpoint \
  -H "Origin: https://example.com" \
  -v
```

### Using Browser DevTools
1. Open your browser's Developer Tools (F12)
2. Navigate to the Network tab
3. Make a request from a different origin
4. Check the response headers for:
   - `Access-Control-Allow-Origin`
   - `Access-Control-Allow-Methods`
   - `Access-Control-Allow-Headers`
   - `Access-Control-Max-Age`

## Troubleshooting

### CORS Error: "Origin not allowed"
**Cause**: The origin is not in the `allowed_origins` list
**Solution**: Add your origin to the configuration or use wildcard for development

### CORS Error: "Method not allowed"
**Cause**: The HTTP method is not in the `allowed_methods` list
**Solution**: Add the required method to the configuration

### CORS Error: "Header not allowed"
**Cause**: A header you're sending is not in the `allowed_headers` list
**Solution**: Add the required header to the configuration

### No CORS headers in response
**Cause**: CORS is disabled or no `Origin` header in request
**Solution**: Enable CORS in configuration and ensure requests include the `Origin` header

## References
- [MDN Web Docs: CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [gin-contrib/cors middleware](https://github.com/gin-contrib/cors)
