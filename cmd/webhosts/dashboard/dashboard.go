package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/bgrewell/dtac-agent/pkg/webhost"
)

//go:embed static
var staticFiles embed.FS

// DashboardModule implements the webhost.WebhostModule interface.
// It serves a simple dashboard with embedded static files and a health endpoint.
type DashboardModule struct {
	config *webhost.InitConfig
	server *http.Server
	logger webhost.Logger
}

// NewDashboardModule creates a new dashboard module instance.
func NewDashboardModule() *DashboardModule {
	return &DashboardModule{}
}

// Name returns the module's unique name.
func (m *DashboardModule) Name() string {
	return "dashboard"
}

// OnInit is called after the host completes initialization.
func (m *DashboardModule) OnInit(ctx context.Context, config *webhost.InitConfig) error {
	m.config = config

	// Extract logger from context
	if logger, ok := ctx.Value("logger").(webhost.Logger); ok && logger != nil {
		m.logger = logger
	} else {
		// Fallback to a basic logger if not provided in context
		m.logger = &basicLogger{}
	}

	m.logger.Info("Dashboard module initialized", map[string]string{
		"module": m.Name(),
		"addr":   config.Listen.Addr,
		"port":   fmt.Sprintf("%d", config.Listen.Port),
	})

	return nil
}

// Run starts the module's HTTP server.
func (m *DashboardModule) Run(ctx context.Context) error {
	// Create HTTP mux
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/healthz", m.healthHandler)

	// Serve embedded static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to create static file system: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// Set up proxies if configured
	if len(m.config.Proxies) > 0 {
		// Get token provider from context (if available)
		var tokenProvider webhost.TokenProvider
		if tp, ok := ctx.Value("tokenProvider").(webhost.TokenProvider); ok {
			tokenProvider = tp
		}
		webhost.SetupProxies(mux, m.config.Proxies, tokenProvider, m.logger)
	}

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", m.config.Listen.Addr, m.config.Listen.Port)
	m.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	m.logger.Info("Starting HTTP server", map[string]string{
		"addr": addr,
	})

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if m.config.Listen.TLSCertPath != "" && m.config.Listen.TLSKeyPath != "" {
			errChan <- m.server.ListenAndServeTLS(m.config.Listen.TLSCertPath, m.config.Listen.TLSKeyPath)
		} else {
			errChan <- m.server.ListenAndServe()
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		m.logger.Info("Shutting down HTTP server", nil)
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := m.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		return nil
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	}
}

// healthHandler handles health check requests.
func (m *DashboardModule) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","module":"%s","timestamp":"%s"}`, m.Name(), time.Now().Format(time.RFC3339))
}

// basicLogger is a fallback logger that writes to stdout.
type basicLogger struct{}

func (l *basicLogger) Log(level webhost.LogLevel, message string, fields map[string]string) {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] [%s] %s %v\n", timestamp, level, message, fields)
}

func (l *basicLogger) Debug(message string, fields map[string]string) {
	l.Log(webhost.LogLevelDebug, message, fields)
}

func (l *basicLogger) Info(message string, fields map[string]string) {
	l.Log(webhost.LogLevelInfo, message, fields)
}

func (l *basicLogger) Warn(message string, fields map[string]string) {
	l.Log(webhost.LogLevelWarn, message, fields)
}

func (l *basicLogger) Error(message string, fields map[string]string) {
	l.Log(webhost.LogLevelError, message, fields)
}

func (l *basicLogger) Fatal(message string, fields map[string]string) {
	l.Log(webhost.LogLevelFatal, message, fields)
}
