package modules

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"io"
	"io/fs"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// WebModuleConfig represents configuration specific to web modules
type WebModuleConfig struct {
	// Port to listen on (0 for auto-assign)
	Port int `json:"port"`
	// StaticPath is the root path for serving static files
	StaticPath string `json:"static_path"`
	// ProxyRoutes defines backend proxy configurations
	ProxyRoutes []ProxyRouteConfig `json:"proxy_routes"`
	// Debug enables request logging
	Debug bool `json:"debug"`
	// RuntimeEnv contains environment variables to expose to the frontend via /config.js
	// These can be used by React/Vite applications that need runtime configuration
	// Example: {"API_URL": "https://api.example.com", "FEATURE_FLAG": "true"}
	RuntimeEnv map[string]string `json:"runtime_env"`
	// InjectConfig automatically injects the config.js script tag into HTML pages
	// When true (default), the module intercepts HTML responses and adds
	// <script src="/config.js"></script> before </head>
	// Set to false to disable automatic injection and manually include the script
	InjectConfig *bool `json:"inject_config"`
}

// ProxyRouteConfig defines a proxy route to a backend service
type ProxyRouteConfig struct {
	// Name is the identifier for this proxy endpoint (e.g., "maas")
	// Requests to /api/<name>/* will be proxied to Target
	Name string `json:"name"`
	// Path is the frontend path to match (e.g., "/api")
	// If not specified and Name is provided, defaults to "/api/<name>"
	Path string `json:"path"`
	// Target is the backend URL to proxy to (e.g., "http://internal.ip/base")
	Target string `json:"target"`
	// StripPath indicates whether to remove the path prefix when proxying
	StripPath bool `json:"strip_path"`
	// AuthType specifies the authentication type (e.g., "bearer", "basic", "none")
	AuthType string `json:"auth_type"`
	// Credentials holds authentication credentials
	Credentials ProxyCredentials `json:"credentials"`
}

// ProxyCredentials holds authentication information for proxy routes
type ProxyCredentials struct {
	// Token for bearer authentication
	Token string `json:"token"`
	// Username for basic authentication
	Username string `json:"username"`
	// Password for basic authentication
	Password string `json:"password"`
	// OAuthConsumerKey for OAuth 1.0a authentication (part 1 of 3-part key)
	OAuthConsumerKey string `json:"oauth_consumer_key"`
	// OAuthToken for OAuth 1.0a authentication (part 2 of 3-part key)
	OAuthToken string `json:"oauth_token"`
	// OAuthTokenSecret for OAuth 1.0a authentication (part 3 of 3-part key)
	OAuthTokenSecret string `json:"oauth_token_secret"`
	// Headers contains custom headers to add to proxied requests
	Headers map[string]string `json:"headers"`
}

// WebModule is a specialized module for hosting web frontends
type WebModule interface {
	Module
	// GetStaticFiles returns the embedded filesystem for static assets
	GetStaticFiles() fs.FS
	// Start starts the web server
	Start() error
	// Stop stops the web server
	Stop() error
	// GetPort returns the port the web server is listening on
	GetPort() int
}

// WebModuleBase provides base functionality for web modules
type WebModuleBase struct {
	ModuleBase
	config       WebModuleConfig
	server       *http.Server
	serverPort   int
	staticFS     embed.FS
	mu           sync.RWMutex
	isRunning    bool
	staticGetter func() fs.FS // Function to get static files from concrete implementation
	httpClient   *http.Client  // Shared HTTP client for proxy requests
}

// Register registers the web module with the module manager
func (w *WebModuleBase) Register(request *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error {
	*reply = api.ModuleRegisterResponse{
		ModuleType:   "web",
		Capabilities: []string{"static_files", "http_server", "logging"},
	}

	// Set default configuration if not provided
	w.config = WebModuleConfig{
		Port:        0, // auto-assign
		StaticPath:  "/",
		ProxyRoutes: []ProxyRouteConfig{},
		Debug:       false, // debug logging disabled by default
	}

	// Log registration
	w.Log(LoggingLevelInfo, "web module registered", map[string]string{
		"module_type": "web",
		"port":        fmt.Sprintf("%d", w.config.Port),
	})

	return nil
}

// GetStaticFiles returns the embedded filesystem - must be overridden by concrete implementations
func (w *WebModuleBase) GetStaticFiles() fs.FS {
	return nil
}

// Start starts the web server
func (w *WebModuleBase) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isRunning {
		return fmt.Errorf("web server is already running")
	}

	// Initialize HTTP client for proxy requests if not already done
	if w.httpClient == nil {
		w.httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	// Get static files using the registered getter function
	var staticFS fs.FS
	if w.staticGetter != nil {
		staticFS = w.staticGetter()
	}

	// Create HTTP server
	mux := http.NewServeMux()

	// Serve /config.js for runtime environment variables
	// This must be registered before static files to take precedence
	mux.HandleFunc("/config.js", w.configJSHandler)

	// Serve static files if filesystem is provided
	if staticFS != nil {
		fileServer := http.FileServer(http.FS(staticFS))
		// Handle root path specially - don't strip prefix for "/"
		if w.config.StaticPath == "/" {
			mux.Handle("/", fileServer)
		} else {
			mux.Handle(w.config.StaticPath, http.StripPrefix(w.config.StaticPath, fileServer))
		}
	}

	// Setup proxy routes
	for _, route := range w.config.ProxyRoutes {
		// Determine the path to handle
		path := route.Path
		if path == "" && route.Name != "" {
			// Default path to /api/<name>
			path = "/api/" + route.Name + "/"
		}
		
		// Ensure path ends with "/" for prefix matching
		if !strings.HasSuffix(path, "/") {
			path = path + "/"
		}

		w.Log(LoggingLevelInfo, "proxy route configured", map[string]string{
			"name":   route.Name,
			"path":   path,
			"target": route.Target,
		})

		// Create proxy handler for this route
		handler := w.createProxyHandler(route, path)
		mux.Handle(path, handler)
	}

	// Determine port
	port := w.config.Port
	if port == 0 {
		port = 8080 // default port
	}

	// Build handler chain
	var handler http.Handler = mux
	
	// Apply HTML injection middleware if enabled (default: true when runtime_env is configured)
	// Note: We check config directly here since we already hold the lock
	injectEnabled := w.config.InjectConfig == nil || *w.config.InjectConfig
	if injectEnabled && len(w.config.RuntimeEnv) > 0 {
		handler = w.htmlInjectionMiddleware(handler)
	}
	
	// Wrap handler with logging middleware
	handler = w.loggingMiddleware(handler)

	w.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	w.serverPort = port
	w.isRunning = true

	// Log static files configuration
	if staticFS != nil {
		w.logStaticFiles(staticFS)
	}

	// Start server in goroutine
	go func() {
		w.Log(LoggingLevelInfo, "starting web server", map[string]string{
			"port":        fmt.Sprintf("%d", w.serverPort),
			"static_path": w.config.StaticPath,
			"debug":       fmt.Sprintf("%t", w.config.Debug),
		})
		if err := w.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			w.Log(LoggingLevelError, "web server error", map[string]string{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

// Stop stops the web server
func (w *WebModuleBase) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isRunning {
		return fmt.Errorf("web server is not running")
	}

	if w.server != nil {
		w.Log(LoggingLevelInfo, "stopping web server", nil)
		if err := w.server.Close(); err != nil {
			return err
		}
	}

	w.isRunning = false
	return nil
}

// GetPort returns the port the web server is listening on
func (w *WebModuleBase) GetPort() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.serverPort
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.statusCode = http.StatusOK
		rw.written = true
	}
	return rw.ResponseWriter.Write(b)
}

// loggingMiddleware wraps an http.Handler and logs request/response details when debug is enabled
func (w *WebModuleBase) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if w.config.Debug {
			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: rw,
				statusCode:     http.StatusOK,
				written:        false,
			}
			
			// Log request
			w.Log(LoggingLevelDebug, "HTTP request received", map[string]string{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			})
			
			// Process request
			next.ServeHTTP(wrapped, r)
			
			// Log response
			w.Log(LoggingLevelDebug, "HTTP response sent", map[string]string{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": fmt.Sprintf("%d", wrapped.statusCode),
				"status_text": http.StatusText(wrapped.statusCode),
			})
		} else {
			next.ServeHTTP(rw, r)
		}
	})
}

// htmlInjectionMiddleware intercepts HTML responses and injects the config.js script tag
func (w *WebModuleBase) htmlInjectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Skip injection for the config.js endpoint itself
		if r.URL.Path == "/config.js" {
			next.ServeHTTP(rw, r)
			return
		}

		// Create a response wrapper to capture the response
		wrapper := &htmlInjectionResponseWriter{
			ResponseWriter: rw,
			request:        r,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// If we buffered HTML content, inject and write it now
		if wrapper.bufferedHTML != nil {
			// Inject the config.js script tag before </head>
			html := wrapper.bufferedHTML
			injectedHTML := injectConfigScript(html)
			
			// Update Content-Length header (must be set before WriteHeader)
			rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(injectedHTML)))
			rw.WriteHeader(wrapper.statusCode)
			rw.Write(injectedHTML)
		}
	})
}

// htmlInjectionResponseWriter wraps http.ResponseWriter to intercept HTML responses
type htmlInjectionResponseWriter struct {
	http.ResponseWriter
	request      *http.Request
	statusCode   int
	bufferedHTML []byte
	wroteHeader  bool
	isHTML       bool
	checked      bool
}

func (w *htmlInjectionResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	
	// Check if this is an HTML response
	if !w.checked {
		w.checked = true
		contentType := w.Header().Get("Content-Type")
		w.isHTML = strings.Contains(contentType, "text/html")
	}
	
	// For non-HTML responses, write header immediately
	if !w.isHTML {
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(code)
	}
	// For HTML responses, delay writing header until we inject the script
}

func (w *htmlInjectionResponseWriter) Write(b []byte) (int, error) {
	// Check content type on first write if not already checked
	if !w.checked {
		w.checked = true
		contentType := w.Header().Get("Content-Type")
		w.isHTML = strings.Contains(contentType, "text/html")
		
		// For non-HTML responses, write header if needed
		if !w.isHTML && !w.wroteHeader {
			w.wroteHeader = true
			w.ResponseWriter.WriteHeader(w.statusCode)
		}
	}

	// For HTML responses, buffer the content
	if w.isHTML {
		w.bufferedHTML = append(w.bufferedHTML, b...)
		return len(b), nil
	}

	// For non-HTML responses, write directly
	if !w.wroteHeader {
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(w.statusCode)
	}
	return w.ResponseWriter.Write(b)
}

// injectConfigScript injects the config.js script tag into HTML content
func injectConfigScript(html []byte) []byte {
	htmlStr := string(html)
	scriptTag := `<script src="/config.js"></script>`
	
	// Try to inject before </head>
	headCloseIdx := strings.Index(strings.ToLower(htmlStr), "</head>")
	if headCloseIdx != -1 {
		return []byte(htmlStr[:headCloseIdx] + scriptTag + "\n" + htmlStr[headCloseIdx:])
	}
	
	// Fallback: inject after <head> if </head> not found
	headOpenIdx := strings.Index(strings.ToLower(htmlStr), "<head>")
	if headOpenIdx != -1 {
		insertPos := headOpenIdx + 6 // length of "<head>"
		return []byte(htmlStr[:insertPos] + "\n" + scriptTag + htmlStr[insertPos:])
	}
	
	// Fallback: inject at the beginning of <body>
	lowerHTML := strings.ToLower(htmlStr)
	bodyOpenIdx := strings.Index(lowerHTML, "<body")
	if bodyOpenIdx != -1 {
		// Find the closing > of the body tag (use lowercase version for consistency)
		closeIdx := strings.Index(lowerHTML[bodyOpenIdx:], ">")
		if closeIdx != -1 {
			insertPos := bodyOpenIdx + closeIdx + 1
			return []byte(htmlStr[:insertPos] + "\n" + scriptTag + htmlStr[insertPos:])
		}
	}
	
	// Last resort: prepend to the content
	return []byte(scriptTag + "\n" + htmlStr)
}

// SetConfig sets the web module configuration
func (w *WebModuleBase) SetConfig(config WebModuleConfig) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.config = config
}

// SetStaticFilesGetter sets the function to retrieve static files from the concrete implementation
func (w *WebModuleBase) SetStaticFilesGetter(getter func() fs.FS) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.staticGetter = getter
}

// configJSHandler serves runtime environment variables as JavaScript
// This allows React/Vite applications to access configuration at runtime
// The output format is: window.__DTAC_CONFIG__ = {...};
func (w *WebModuleBase) configJSHandler(rw http.ResponseWriter, r *http.Request) {
	w.mu.RLock()
	runtimeEnv := w.config.RuntimeEnv
	w.mu.RUnlock()

	// Use json.Marshal for proper escaping of all values
	// This handles all special characters including <, >, &, quotes, newlines, etc.
	configJSON, err := json.Marshal(runtimeEnv)
	if err != nil {
		w.Log(LoggingLevelError, "failed to marshal runtime config", map[string]string{
			"error": err.Error(),
		})
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Build the JavaScript output with proper escaping
	// json.Marshal escapes <, >, & to unicode escape sequences by default
	var sb strings.Builder
	sb.WriteString("// DTAC Runtime Configuration\n")
	sb.WriteString("// Generated by DTAC web module - do not edit\n")
	sb.WriteString("window.__DTAC_CONFIG__ = ")
	
	if runtimeEnv == nil {
		sb.WriteString("{}")
	} else {
		sb.Write(configJSON)
	}
	sb.WriteString(";\n")

	// Set content type and cache control headers
	rw.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Expires", "0")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(sb.String()))
}

// logStaticFiles logs information about available static files and routes
func (w *WebModuleBase) logStaticFiles(staticFS fs.FS) {
	fileCount := 0
	var files []string
	
	// Walk the filesystem to enumerate files
	fs.WalkDir(staticFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fileCount++
			// Build the URL path
			urlPath := w.config.StaticPath
			if urlPath == "/" {
				urlPath = "/" + path
			} else {
				urlPath = urlPath + "/" + path
			}
			files = append(files, urlPath)
		}
		return nil
	})
	
	// Log summary
	w.Log(LoggingLevelInfo, "static files configured", map[string]string{
		"file_count":  fmt.Sprintf("%d", fileCount),
		"static_path": w.config.StaticPath,
	})
	
	// Log individual file routes if in debug mode
	if w.config.Debug {
		for _, filePath := range files {
			w.Log(LoggingLevelDebug, "static file route", map[string]string{
				"route": fmt.Sprintf("http://localhost:%d%s", w.serverPort, filePath),
				"path":  filePath,
			})
		}
	}
}

// normalizePath ensures a path starts with "/" if not empty, or returns "/" if empty
func normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

// createProxyHandler creates an HTTP handler that proxies requests to a backend service
func (w *WebModuleBase) createProxyHandler(route ProxyRouteConfig, frontendPath string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Parse target URL
		targetURL, err := url.Parse(route.Target)
		if err != nil {
			w.Log(LoggingLevelError, "invalid proxy target URL", map[string]string{
				"target": route.Target,
				"error":  err.Error(),
			})
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Build the backend path
		backendPath := r.URL.Path
		if route.StripPath {
			// Remove the frontend path prefix
			prefix := strings.TrimSuffix(frontendPath, "/")
			if strings.HasPrefix(backendPath, prefix) {
				backendPath = strings.TrimPrefix(backendPath, prefix)
			}
		}

		// Normalize the backend path
		backendPath = normalizePath(backendPath)

		// Combine target base path with backend path
		targetURL.Path = strings.TrimSuffix(targetURL.Path, "/") + backendPath
		targetURL.RawQuery = r.URL.RawQuery

		// Create the proxy request
		proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
		if err != nil {
			w.Log(LoggingLevelError, "failed to create proxy request", map[string]string{
				"error": err.Error(),
			})
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Copy headers from original request
		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Add authentication headers based on auth type
		w.addAuthHeaders(proxyReq, route)

		// Add custom headers from credentials
		for key, value := range route.Credentials.Headers {
			proxyReq.Header.Set(key, value)
		}

		// Log the proxy request if debug is enabled
		if w.config.Debug {
			w.Log(LoggingLevelDebug, "proxying request", map[string]string{
				"method":       r.Method,
				"from":         r.URL.Path,
				"to":           targetURL.String(),
				"auth_type":    route.AuthType,
				"backend_path": backendPath,
			})
		}

		// Execute the proxy request using shared HTTP client
		resp, err := w.httpClient.Do(proxyReq)
		if err != nil {
			w.Log(LoggingLevelError, "proxy request failed", map[string]string{
				"target": targetURL.String(),
				"error":  err.Error(),
			})
			http.Error(rw, "Bad Gateway", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(key, value)
			}
		}

		// Write status code
		rw.WriteHeader(resp.StatusCode)

		// Copy response body
		_, err = io.Copy(rw, resp.Body)
		if err != nil {
			w.Log(LoggingLevelError, "failed to copy proxy response", map[string]string{
				"error": err.Error(),
			})
		}

		// Log response if debug is enabled
		if w.config.Debug {
			w.Log(LoggingLevelDebug, "proxy response", map[string]string{
				"status_code": fmt.Sprintf("%d", resp.StatusCode),
				"status_text": resp.Status,
			})
		}
	})
}

// addAuthHeaders adds authentication headers to the proxy request based on the auth type
func (w *WebModuleBase) addAuthHeaders(req *http.Request, route ProxyRouteConfig) {
	authType := strings.ToLower(route.AuthType)
	
	switch authType {
	case "bearer":
		if route.Credentials.Token != "" {
			req.Header.Set("Authorization", "Bearer "+route.Credentials.Token)
		}
	case "basic":
		if route.Credentials.Username != "" {
			req.SetBasicAuth(route.Credentials.Username, route.Credentials.Password)
		}
	case "oauth":
		if route.Credentials.OAuthConsumerKey != "" {
			// Generate OAuth 1.0a header (PLAINTEXT signature method as used by MAAS)
			oauthHeader := w.generateOAuthHeader(
				route.Credentials.OAuthConsumerKey,
				route.Credentials.OAuthToken,
				route.Credentials.OAuthTokenSecret,
			)
			req.Header.Set("Authorization", oauthHeader)
		}
	case "none", "":
		// No authentication
	default:
		w.Log(LoggingLevelWarning, "unknown auth type", map[string]string{
			"auth_type": route.AuthType,
		})
	}
}

// generateOAuthHeader generates an OAuth 1.0a authorization header using PLAINTEXT signature method
// This is compatible with MAAS API authentication
func (w *WebModuleBase) generateOAuthHeader(consumerKey, token, tokenSecret string) string {
	// Generate random nonce
	nonce := generateNonce()
	
	// Get current timestamp
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	
	// Build OAuth header using PLAINTEXT signature method
	// Signature format: &{token_secret} (note the & prefix and no consumer secret)
	signature := fmt.Sprintf("%%26%s", tokenSecret)
	
	header := fmt.Sprintf(
		"OAuth oauth_consumer_key=\"%s\",oauth_token=\"%s\",oauth_signature_method=\"PLAINTEXT\","+
			"oauth_timestamp=\"%s\",oauth_nonce=\"%s\",oauth_version=\"1.0\",oauth_signature=\"%s\"",
		consumerKey, token, timestamp, nonce, signature,
	)
	
	return header
}

// generateNonce generates a random nonce for OAuth requests
func generateNonce() string {
	// Generate a random number up to 10^10
	max := big.NewInt(10000000000) // 10^10
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback to timestamp-based nonce if crypto/rand fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return n.String()
}

// ParseWebModuleConfig parses a configuration map into a WebModuleConfig struct
// This centralizes all configuration parsing logic for web modules
func ParseWebModuleConfig(configMap map[string]interface{}) WebModuleConfig {
	config := WebModuleConfig{
		Port:        8080, // default port
		StaticPath:  "/",
		ProxyRoutes: []ProxyRouteConfig{},
		Debug:       false,
	}
	
	// Parse port
	if port, ok := configMap["port"]; ok {
		if portFloat, ok := port.(float64); ok {
			config.Port = int(portFloat)
		} else if portInt, ok := port.(int); ok {
			config.Port = portInt
		}
	}
	
	// Parse static_path
	if staticPath, ok := configMap["static_path"].(string); ok {
		config.StaticPath = staticPath
	}
	
	// Parse debug
	if debug, ok := configMap["debug"].(bool); ok {
		config.Debug = debug
	}
	
	// Parse proxy_routes
	if proxyRoutes, ok := configMap["proxy_routes"]; ok {
		if routesSlice, ok := proxyRoutes.([]interface{}); ok {
			for _, routeInterface := range routesSlice {
				if routeMap, ok := routeInterface.(map[string]interface{}); ok {
					route := parseProxyRoute(routeMap)
					config.ProxyRoutes = append(config.ProxyRoutes, route)
				}
			}
		}
	}

	// Parse runtime_env for frontend environment variables
	if runtimeEnv, ok := configMap["runtime_env"]; ok {
		if envMap, ok := runtimeEnv.(map[string]interface{}); ok {
			config.RuntimeEnv = make(map[string]string)
			for k, v := range envMap {
				if strVal, ok := v.(string); ok {
					config.RuntimeEnv[k] = strVal
				} else {
					// Convert non-string values to string
					config.RuntimeEnv[k] = fmt.Sprintf("%v", v)
				}
			}
		}
	}

	// Parse inject_config (defaults to true if not specified)
	if injectConfig, ok := configMap["inject_config"].(bool); ok {
		config.InjectConfig = &injectConfig
	}
	
	return config
}

// parseProxyRoute extracts a single proxy route configuration from a map
func parseProxyRoute(routeMap map[string]interface{}) ProxyRouteConfig {
	route := ProxyRouteConfig{}
	
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
func parseProxyCredentials(credsMap map[string]interface{}) ProxyCredentials {
	creds := ProxyCredentials{}
	
	// Parse token (for bearer auth)
	if token, ok := credsMap["token"].(string); ok {
		creds.Token = token
	}
	
	// Parse username (for basic auth)
	if username, ok := credsMap["username"].(string); ok {
		creds.Username = username
	}
	
	// Parse password (for basic auth)
	if password, ok := credsMap["password"].(string); ok {
		creds.Password = password
	}
	
	// Parse OAuth credentials
	if oauthConsumerKey, ok := credsMap["oauth_consumer_key"].(string); ok {
		creds.OAuthConsumerKey = oauthConsumerKey
	}
	
	if oauthToken, ok := credsMap["oauth_token"].(string); ok {
		creds.OAuthToken = oauthToken
	}
	
	if oauthTokenSecret, ok := credsMap["oauth_token_secret"].(string); ok {
		creds.OAuthTokenSecret = oauthTokenSecret
	}
	
	// Parse headers (for custom header auth)
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
