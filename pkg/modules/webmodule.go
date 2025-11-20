package modules

import (
	"embed"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"io/fs"
	"net/http"
	"sync"
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
}

// ProxyRouteConfig defines a proxy route to a backend service
type ProxyRouteConfig struct {
	// Path is the frontend path to match (e.g., "/api")
	Path string `json:"path"`
	// Target is the backend URL to proxy to
	Target string `json:"target"`
	// StripPath indicates whether to remove the path prefix when proxying
	StripPath bool `json:"strip_path"`
	// InjectToken indicates whether to inject auth tokens
	InjectToken bool `json:"inject_token"`
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
	config     WebModuleConfig
	server     *http.Server
	serverPort int
	staticFS   embed.FS
	mu         sync.RWMutex
	isRunning  bool
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

	// Create HTTP server
	mux := http.NewServeMux()

	// Serve static files if filesystem is provided
	if staticFS := w.GetStaticFiles(); staticFS != nil {
		fileServer := http.FileServer(http.FS(staticFS))
		// Handle root path specially - don't strip prefix for "/"
		if w.config.StaticPath == "/" {
			mux.Handle("/", fileServer)
		} else {
			mux.Handle(w.config.StaticPath, http.StripPrefix(w.config.StaticPath, fileServer))
		}
	}

	// Setup proxy routes (placeholder - will be implemented in future iterations)
	for _, route := range w.config.ProxyRoutes {
		w.Log(LoggingLevelDebug, "proxy route configured", map[string]string{
			"path":   route.Path,
			"target": route.Target,
		})
		// TODO: Implement proxy handler
	}

	// Determine port
	port := w.config.Port
	if port == 0 {
		port = 8080 // default port
	}

	// Wrap handler with logging middleware
	handler := w.loggingMiddleware(mux)

	w.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	w.serverPort = port
	w.isRunning = true

	// Log static files configuration
	if staticFS := w.GetStaticFiles(); staticFS != nil {
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

// SetConfig sets the web module configuration
func (w *WebModuleBase) SetConfig(config WebModuleConfig) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.config = config
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
