package webhost

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	// EnvDTACWebhosts is the environment variable that indicates managed mode.
	EnvDTACWebhosts = "DTAC_WEBHOSTS"

	// Context keys for logger and token provider
	contextKeyLogger        = "logger"
	contextKeyTokenProvider = "tokenProvider"
)

// WebhostHost manages the lifecycle of a webhost module.
// It handles mode detection, initialization, and coordination between
// the module and DTAC (in managed mode) or local config (in standalone mode).
type WebhostHost interface {
	// Run starts the host and blocks until the context is cancelled or an error occurs.
	Run(ctx context.Context) error

	// GetMode returns the current host mode (standalone or managed).
	GetMode() HostMode

	// GetLogger returns the logger instance.
	GetLogger() Logger

	// GetTokenProvider returns the token provider instance.
	GetTokenProvider() TokenProvider
}

// DefaultWebhostHost is the standard implementation of WebhostHost.
type DefaultWebhostHost struct {
	module        WebhostModule
	mode          HostMode
	logger        Logger
	tokenProvider TokenProvider
	grpcPort      int // gRPC port for managed mode
}

// NewDefaultWebhostHost creates a new DefaultWebhostHost for the given module.
// It automatically detects the mode based on environment variables.
func NewDefaultWebhostHost(module WebhostModule) *DefaultWebhostHost {
	mode := detectMode()
	return &DefaultWebhostHost{
		module: module,
		mode:   mode,
	}
}

// detectMode checks environment variables to determine if we're in managed mode.
func detectMode() HostMode {
	if os.Getenv(EnvDTACWebhosts) == "true" {
		return HostModeManaged
	}
	return HostModeStandalone
}

// Run starts the host and executes the module.
func (h *DefaultWebhostHost) Run(ctx context.Context) error {
	if h.mode == HostModeManaged {
		return h.runManaged(ctx)
	}
	return h.runStandalone(ctx)
}

// GetMode returns the current host mode.
func (h *DefaultWebhostHost) GetMode() HostMode {
	return h.mode
}

// GetLogger returns the logger instance.
func (h *DefaultWebhostHost) GetLogger() Logger {
	return h.logger
}

// GetTokenProvider returns the token provider instance.
func (h *DefaultWebhostHost) GetTokenProvider() TokenProvider {
	return h.tokenProvider
}

// runManaged runs the module in DTAC-managed mode.
// This involves:
// 1. Starting a gRPC server that implements WebhostService
// 2. Printing a CONNECT handshake line for DTAC to discover the gRPC endpoint
// 3. Waiting for DTAC to call Init
// 4. Running the module's main logic
func (h *DefaultWebhostHost) runManaged(ctx context.Context) error {
	// TODO: Implement full gRPC server setup
	// For now, this is a skeleton that demonstrates the flow

	// Initialize managed-mode logger (sends to EventStream)
	h.logger = newManagedLogger()

	// Initialize managed-mode token provider (requests from DTAC)
	h.tokenProvider = newManagedTokenProvider()

	// Start gRPC server on a random port
	// In the full implementation, this would:
	// - Create a gRPC server
	// - Register WebhostService implementation
	// - Listen on a port (OS-assigned or configured)
	// - Print CONNECT{{...}} handshake
	h.grpcPort = 50051 // Placeholder

	// Print CONNECT handshake for DTAC to discover us
	connectInfo := map[string]interface{}{
		"protocol": "grpc",
		"address":  fmt.Sprintf("127.0.0.1:%d", h.grpcPort),
		"version":  "1.0",
	}
	connectJSON, _ := json.Marshal(connectInfo)
	fmt.Printf("CONNECT{{%s}}\n", connectJSON)

	h.logger.Info("Webhost module started in managed mode", map[string]string{
		"module": h.module.Name(),
		"mode":   string(h.mode),
	})

	// In the full implementation, we'd wait for DTAC to call Init via gRPC
	// For this skeleton, we'll create a minimal config and call OnInit
	config := &InitConfig{
		ModuleName: h.module.Name(),
		Listen: ListenConfig{
			Addr: "127.0.0.1",
			Port: 8080,
		},
		Proxies:     []ProxyConfig{},
		Permissions: Permissions{},
		Config:      make(map[string]interface{}),
	}

	// Create a new context with logger and token provider
	ctx = context.WithValue(ctx, contextKeyLogger, h.logger)
	ctx = context.WithValue(ctx, contextKeyTokenProvider, h.tokenProvider)

	if err := h.module.OnInit(ctx, config); err != nil {
		return fmt.Errorf("module OnInit failed: %w", err)
	}

	// Run the module
	return h.module.Run(ctx)
}

// runStandalone runs the module in standalone mode.
// This involves:
// 1. Reading configuration from CLI flags or config file
// 2. Creating a synthetic InitConfig
// 3. Initializing stdio logger
// 4. Running the module's main logic
func (h *DefaultWebhostHost) runStandalone(ctx context.Context) error {
	// Initialize standalone logger (writes to stdout/stderr)
	h.logger = newStdioLogger()

	// Initialize standalone token provider (returns errors or mock tokens)
	h.tokenProvider = newStandaloneTokenProvider()

	h.logger.Info("Webhost module started in standalone mode", map[string]string{
		"module": h.module.Name(),
		"mode":   string(h.mode),
	})

	// In a real implementation, we'd read from flags or config file
	// For this skeleton, we'll use hardcoded defaults
	config := &InitConfig{
		ModuleName: h.module.Name(),
		Listen: ListenConfig{
			Addr: "127.0.0.1",
			Port: 8080,
		},
		Proxies:     []ProxyConfig{},
		Permissions: Permissions{},
		Config:      make(map[string]interface{}),
	}

	// Create a new context with logger and token provider
	ctx = context.WithValue(ctx, contextKeyLogger, h.logger)
	ctx = context.WithValue(ctx, contextKeyTokenProvider, h.tokenProvider)

	if err := h.module.OnInit(ctx, config); err != nil {
		return fmt.Errorf("module OnInit failed: %w", err)
	}

	// Run the module
	return h.module.Run(ctx)
}

// managedLogger implements Logger for managed mode (sends to EventStream).
type managedLogger struct{}

func newManagedLogger() Logger {
	return &managedLogger{}
}

func (l *managedLogger) Log(level LogLevel, message string, fields map[string]string) {
	// TODO: Send to EventStream via gRPC
	// For now, fallback to stdout
	fmt.Printf("[%s] %s %v\n", level, message, fields)
}

func (l *managedLogger) Debug(message string, fields map[string]string) {
	l.Log(LogLevelDebug, message, fields)
}

func (l *managedLogger) Info(message string, fields map[string]string) {
	l.Log(LogLevelInfo, message, fields)
}

func (l *managedLogger) Warn(message string, fields map[string]string) {
	l.Log(LogLevelWarn, message, fields)
}

func (l *managedLogger) Error(message string, fields map[string]string) {
	l.Log(LogLevelError, message, fields)
}

func (l *managedLogger) Fatal(message string, fields map[string]string) {
	l.Log(LogLevelFatal, message, fields)
	os.Exit(1)
}

// stdioLogger implements Logger for standalone mode (writes to stdout/stderr).
type stdioLogger struct{}

func newStdioLogger() Logger {
	return &stdioLogger{}
}

func (l *stdioLogger) Log(level LogLevel, message string, fields map[string]string) {
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] [%s] %s %v\n", timestamp, level, message, fields)
}

func (l *stdioLogger) Debug(message string, fields map[string]string) {
	l.Log(LogLevelDebug, message, fields)
}

func (l *stdioLogger) Info(message string, fields map[string]string) {
	l.Log(LogLevelInfo, message, fields)
}

func (l *stdioLogger) Warn(message string, fields map[string]string) {
	l.Log(LogLevelWarn, message, fields)
}

func (l *stdioLogger) Error(message string, fields map[string]string) {
	l.Log(LogLevelError, message, fields)
}

func (l *stdioLogger) Fatal(message string, fields map[string]string) {
	l.Log(LogLevelFatal, message, fields)
	os.Exit(1)
}

// managedTokenProvider implements TokenProvider for managed mode.
type managedTokenProvider struct{}

func newManagedTokenProvider() TokenProvider {
	return &managedTokenProvider{}
}

func (p *managedTokenProvider) GetToken(ctx context.Context, scopes []string, forProxy string) (string, time.Time, error) {
	// TODO: Request token from DTAC via gRPC RequestToken RPC
	// For this skeleton, return an error
	return "", time.Time{}, errors.New("token provider not fully implemented")
}

// standaloneTokenProvider implements TokenProvider for standalone mode.
type standaloneTokenProvider struct{}

func newStandaloneTokenProvider() TokenProvider {
	return &standaloneTokenProvider{}
}

func (p *standaloneTokenProvider) GetToken(ctx context.Context, scopes []string, forProxy string) (string, time.Time, error) {
	// In standalone mode, we don't have DTAC to issue tokens
	// Return a mock token or error depending on configuration
	return "", time.Time{}, errors.New("token provider not available in standalone mode")
}
