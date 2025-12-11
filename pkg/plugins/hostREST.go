package plugins

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// RESTPluginHost is the REST server implementation for standalone plugin mode
type RESTPluginHost struct {
	Plugin   Plugin
	Config   *StandaloneConfig
	server   *http.Server
	mux      *http.ServeMux
	port     int
	endpoint []*endpoint.Endpoint
}

// NewRESTPluginHost creates a new REST plugin host
func NewRESTPluginHost(plugin Plugin, config *StandaloneConfig) (*RESTPluginHost, error) {
	host := &RESTPluginHost{
		Plugin: plugin,
		Config: config,
		mux:    http.NewServeMux(),
		port:   config.Port,
	}

	return host, nil
}

// Serve starts the REST server
func (rh *RESTPluginHost) Serve() error {
	// Register the plugin to get its endpoints
	if err := rh.registerPlugin(); err != nil {
		return fmt.Errorf("failed to register plugin: %w", err)
	}

	// Setup routes
	rh.setupRoutes()

	// Configure server address
	addr := fmt.Sprintf("%s:%d", rh.Config.Host, rh.Config.Port)

	// Create HTTP server
	rh.server = &http.Server{
		Addr:         addr,
		Handler:      rh.mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Configure TLS if needed
	if rh.Config.Protocol == "https" && rh.Config.TLSCertPath != "" && rh.Config.TLSKeyPath != "" {
		// Load TLS certificates
		cert, err := tls.LoadX509KeyPair(rh.Config.TLSCertPath, rh.Config.TLSKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificates: %w", err)
		}

		rh.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}

		log.Printf("Starting %s REST server on %s (HTTPS)\n", rh.Plugin.Name(), addr)
		return rh.server.ListenAndServeTLS("", "")
	}

	log.Printf("Starting %s REST server on %s (HTTP)\n", rh.Plugin.Name(), addr)
	return rh.server.ListenAndServe()
}

// GetPort returns the port the server is listening on
func (rh *RESTPluginHost) GetPort() int {
	return rh.port
}

// Shutdown gracefully shuts down the REST server
func (rh *RESTPluginHost) Shutdown(ctx context.Context) error {
	if rh.server != nil {
		return rh.server.Shutdown(ctx)
	}
	return nil
}

// registerPlugin calls the plugin's Register method to get endpoints
func (rh *RESTPluginHost) registerPlugin() error {
	// We need to call the plugin's Register method to initialize its endpoints
	// Create a minimal register request
	request := &api.RegisterRequest{
		Config:        "{}",
		DefaultSecure: false,
	}
	
	response := &api.RegisterResponse{}
	
	err := rh.Plugin.Register(request, response)
	if err != nil {
		return fmt.Errorf("failed to register plugin: %w", err)
	}
	
	return nil
}

// setupRoutes configures HTTP routes for plugin endpoints
func (rh *RESTPluginHost) setupRoutes() {
	// Add health check endpoint
	rh.mux.HandleFunc("/health", rh.handleHealth)
	
	// Add a generic handler that will route to plugin methods
	rh.mux.HandleFunc("/", rh.handlePluginRequest)
}

// handleHealth handles health check requests
func (rh *RESTPluginHost) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"plugin": rh.Plugin.Name(),
	})
}

// handlePluginRequest handles all plugin endpoint requests
func (rh *RESTPluginHost) handlePluginRequest(w http.ResponseWriter, r *http.Request) {
	// Skip health endpoint
	if r.URL.Path == "/health" {
		return
	}

	// Extract method name from path and HTTP method
	// Format: /{rootPath}/{endpoint}
	// HTTP method maps to endpoint.Action
	path := strings.TrimPrefix(r.URL.Path, "/")
	
	// Remove root path prefix if present
	rootPath := rh.Plugin.RootPath()
	if rootPath != "" && strings.HasPrefix(path, rootPath+"/") {
		path = strings.TrimPrefix(path, rootPath+"/")
	} else if rootPath != "" && path == rootPath {
		path = ""
	}

	// Map HTTP method to endpoint action
	var action endpoint.Action
	switch r.Method {
	case http.MethodGet:
		action = endpoint.ActionRead
	case http.MethodPost:
		action = endpoint.ActionCreate
	case http.MethodPut:
		action = endpoint.ActionWrite
	case http.MethodDelete:
		action = endpoint.ActionDelete
	case http.MethodPatch:
		action = endpoint.ActionWrite
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Build method name: "ACTION:path"
	methodName := fmt.Sprintf("%s:%s", action, path)

	// Build endpoint.Request from HTTP request
	req, err := rh.buildEndpointRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to build request: %v", err), http.StatusBadRequest)
		return
	}

	// Call the plugin
	resp, err := rh.Plugin.Call(methodName, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Plugin call failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Send response
	rh.sendEndpointResponse(w, resp)
}

// buildEndpointRequest converts an HTTP request to an endpoint.Request
func (rh *RESTPluginHost) buildEndpointRequest(r *http.Request) (*endpoint.Request, error) {
	req := &endpoint.Request{
		Metadata:   make(map[string]string),
		Headers:    make(map[string][]string),
		Parameters: make(map[string][]string),
	}

	// Copy headers
	for key, values := range r.Header {
		req.Headers[key] = values
	}

	// Copy query parameters
	for key, values := range r.URL.Query() {
		req.Parameters[key] = values
	}

	// Read body
	if r.Body != nil {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		req.Body = body
	}

	return req, nil
}

// sendEndpointResponse sends an endpoint.Response as an HTTP response
func (rh *RESTPluginHost) sendEndpointResponse(w http.ResponseWriter, resp *endpoint.Response) {
	// Set headers
	if resp.Headers != nil {
		for key, values := range resp.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}

	// Default to JSON content type if not set
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	// Write status code (default 200)
	w.WriteHeader(http.StatusOK)

	// Write body
	if resp.Value != nil {
		w.Write(resp.Value)
	}
}
