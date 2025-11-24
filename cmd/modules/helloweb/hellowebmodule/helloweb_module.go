package hellowebmodule

import (
	"embed"
	"encoding/json"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/modules"
	"io/fs"
	"reflect"
)

//go:embed static
var staticFiles embed.FS

// This sets a non-existent variable to the interface type of module then attempts to assign
// a pointer to HelloWebModule to it. This ensures the HelloWebModule struct implements the Module interface.
var _ modules.Module = &HelloWebModule{}

// NewHelloWebModule is a constructor that returns a new instance of the HelloWebModule
func NewHelloWebModule() *HelloWebModule {
	// Create a new instance of the module
	hwm := &HelloWebModule{
		WebModuleBase: modules.WebModuleBase{},
	}
	
	// Set root path
	hwm.SetRootPath("helloweb")

	// Register the static files getter so the base can call our GetStaticFiles method
	hwm.SetStaticFilesGetter(func() fs.FS {
		return hwm.GetStaticFiles()
	})

	// Set default configuration with debug enabled
	hwm.SetConfig(modules.WebModuleConfig{
		Port:        8090,
		StaticPath:  "/",
		ProxyRoutes: []modules.ProxyRouteConfig{},
		Debug:       true, // Enable debug logging by default for this example
	})

	// Return the new instance
	return hwm
}

// HelloWebModule is the module struct that implements the Module interface
type HelloWebModule struct {
	// WebModuleBase provides web server functionality
	modules.WebModuleBase
}

// Name returns the name of the module type
func (h HelloWebModule) Name() string {
	t := reflect.TypeOf(h)
	return t.Name()
}

// GetStaticFiles returns the embedded filesystem for static assets
func (h *HelloWebModule) GetStaticFiles() fs.FS {
	// Return the static subdirectory from the embedded filesystem
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		h.Log(modules.LoggingLevelError, "failed to get static files subdirectory", map[string]string{
			"error": err.Error(),
		})
		return nil
	}
	return staticFS
}

// Register registers the module with the module manager
func (h *HelloWebModule) Register(request *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error {
	*reply = api.ModuleRegisterResponse{
		ModuleType:   "web",
		Capabilities: []string{"static_files", "http_server", "logging"},
		Endpoints:    make([]*api.PluginEndpoint, 0),
	}

	// Parse configuration
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return err
	}

	// Build web config from parsed values
	webConfig := modules.WebModuleConfig{
		Port:        8090, // default
		StaticPath:  "/",
		ProxyRoutes: []modules.ProxyRouteConfig{},
		Debug:       true, // default to debug enabled for examples
	}

	// Update port if provided in config
	if port, ok := config["port"]; ok {
		if portFloat, ok := port.(float64); ok {
			webConfig.Port = int(portFloat)
		}
	}

	// Update debug if provided in config
	if debug, ok := config["debug"]; ok {
		if debugBool, ok := debug.(bool); ok {
			webConfig.Debug = debugBool
		}
	}

	// Parse proxy routes if provided in config
	webConfig.ProxyRoutes = parseProxyRoutes(config)

	h.SetConfig(webConfig)

	// Log registration
	h.Log(modules.LoggingLevelInfo, "hello web module registered", map[string]string{
		"module_type": "web",
		"port":        fmt.Sprintf("%d", h.GetPort()),
	})

	// Start the web server
	err = h.Start()
	if err != nil {
		h.Log(modules.LoggingLevelError, "failed to start web server", map[string]string{
			"error": err.Error(),
		})
		return err
	}

	h.Log(modules.LoggingLevelInfo, "web server started successfully", map[string]string{
		"port": fmt.Sprintf("%d", h.GetPort()),
	})

	// Note: Web modules can optionally expose API endpoints through DTAC
	// For this example, we're only serving static files
	// To add API endpoints, create endpoint.Endpoint objects, register them with RegisterMethods,
	// and add them to reply.Endpoints

	return nil
}

// parseProxyRoutes extracts proxy route configurations from the config map
func parseProxyRoutes(config map[string]interface{}) []modules.ProxyRouteConfig {
	routes := []modules.ProxyRouteConfig{}
	
	proxyRoutes, ok := config["proxy_routes"]
	if !ok {
		return routes
	}
	
	routesSlice, ok := proxyRoutes.([]interface{})
	if !ok {
		return routes
	}
	
	for _, routeInterface := range routesSlice {
		routeMap, ok := routeInterface.(map[string]interface{})
		if !ok {
			continue
		}
		
		route := parseProxyRoute(routeMap)
		routes = append(routes, route)
	}
	
	return routes
}

// parseProxyRoute extracts a single proxy route configuration from a map
func parseProxyRoute(routeMap map[string]interface{}) modules.ProxyRouteConfig {
	route := modules.ProxyRouteConfig{}
	
	// Parse name
	if name, ok := routeMap["name"].(string); ok {
		route.Name = name
	}
	
	// Parse path (optional)
	if path, ok := routeMap["path"].(string); ok {
		route.Path = path
	}
	
	// Parse target
	if target, ok := routeMap["target"].(string); ok {
		route.Target = target
	}
	
	// Parse strip_path
	if stripPath, ok := routeMap["strip_path"].(bool); ok {
		route.StripPath = stripPath
	}
	
	// Parse auth_type
	if authType, ok := routeMap["auth_type"].(string); ok {
		route.AuthType = authType
	}
	
	// Parse credentials
	if credsInterface, ok := routeMap["credentials"]; ok {
		if credsMap, ok := credsInterface.(map[string]interface{}); ok {
			route.Credentials = parseProxyCredentials(credsMap)
		}
	}
	
	return route
}

// parseProxyCredentials extracts credential information from a map
func parseProxyCredentials(credsMap map[string]interface{}) modules.ProxyCredentials {
	creds := modules.ProxyCredentials{}
	
	// Parse token
	if token, ok := credsMap["token"].(string); ok {
		creds.Token = token
	}
	
	// Parse username
	if username, ok := credsMap["username"].(string); ok {
		creds.Username = username
	}
	
	// Parse password
	if password, ok := credsMap["password"].(string); ok {
		creds.Password = password
	}
	
	// Parse headers
	if headersInterface, ok := credsMap["headers"]; ok {
		if headersMap, ok := headersInterface.(map[string]interface{}); ok {
			creds.Headers = make(map[string]string)
			for k, v := range headersMap {
				if strVal, ok := v.(string); ok {
					creds.Headers[k] = strVal
				}
			}
		}
	}
	
	return creds
}
