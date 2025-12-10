package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bgrewell/dtac-agent/cmd/plugins/standalone-hello/standalonehello"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"github.com/bgrewell/dtac-agent/pkg/plugins/standalone"
)

func main() {
	// Parse command-line flags
	standaloneMode := flag.Bool("standalone", false, "Run as standalone REST server")
	port := flag.Int("port", 8080, "Port to listen on (standalone mode only)")
	enableCORS := flag.Bool("cors", true, "Enable CORS (standalone mode only)")
	enableTLS := flag.Bool("tls", false, "Enable TLS/HTTPS (standalone mode only)")
	certFile := flag.String("cert", "", "TLS certificate file (required if -tls is true)")
	keyFile := flag.String("key", "", "TLS key file (required if -tls is true)")
	logLevel := flag.String("log-level", "info", "Log level: debug, info, warn, error")
	flag.Parse()

	// Create the plugin instance
	plugin := standalonehello.NewStandaloneHelloPlugin()

	if *standaloneMode {
		// Run in standalone mode with REST API
		runStandalone(plugin, *port, *enableCORS, *enableTLS, *certFile, *keyFile, *logLevel)
	} else {
		// Run in traditional plugin mode (via gRPC with DTAC agent)
		runTraditional(plugin)
	}
}

// runStandalone runs the plugin as a standalone REST server
func runStandalone(plugin *standalonehello.StandaloneHelloPlugin, port int, enableCORS, enableTLS bool, certFile, keyFile, logLevel string) {
	fmt.Printf("Starting plugin in standalone mode...\n")
	fmt.Printf("Port: %d\n", port)
	fmt.Printf("CORS: %v\n", enableCORS)
	fmt.Printf("TLS: %v\n", enableTLS)
	fmt.Printf("Log Level: %s\n\n", logLevel)

	// Validate TLS configuration
	if enableTLS {
		if certFile == "" || keyFile == "" {
			log.Fatal("Error: -cert and -key are required when -tls is enabled")
		}
	}

	// Create standalone REST adapter configuration
	config := &standalone.StandaloneRESTConfig{
		Port:       port,
		EnableTLS:  enableTLS,
		CertFile:   certFile,
		KeyFile:    keyFile,
		EnableCORS: enableCORS,
		LogLevel:   logLevel,
	}

	// Create the standalone adapter
	adapter, err := standalone.NewStandaloneRESTAdapter(plugin, config)
	if err != nil {
		log.Fatalf("Failed to create standalone adapter: %v", err)
	}

	// Start the server
	ctx := context.Background()
	if err := adapter.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	protocol := "http"
	if enableTLS {
		protocol = "https"
	}

	fmt.Printf("\n✓ Standalone REST server started successfully!\n\n")
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("  GET  %s://localhost:%d/hello/message - Get hello message\n", protocol, port)
	fmt.Printf("  POST %s://localhost:%d/hello/echo    - Echo a message\n\n", protocol, port)
	fmt.Printf("Example curl commands:\n")
	fmt.Printf("  curl %s://localhost:%d/hello/message\n", protocol, port)
	fmt.Printf("  curl -X POST -H 'Content-Type: application/json' -d '{\"message\":\"test\"}' %s://localhost:%d/hello/echo\n\n", protocol, port)
	fmt.Printf("Press Ctrl+C to stop...\n\n")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down gracefully...")
	if err := adapter.Stop(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	fmt.Println("Server stopped")
}

// runTraditional runs the plugin in traditional mode via gRPC
func runTraditional(plugin *standalonehello.StandaloneHelloPlugin) {
	fmt.Println("Starting plugin in traditional mode (gRPC)...")

	// Create the plugin host
	host, err := plugins.NewPluginHost(plugin)
	if err != nil {
		log.Fatalf("Failed to create plugin host: %v", err)
	}

	// Serve the plugin (this will block and communicate with DTAC agent)
	if err := host.Serve(); err != nil {
		log.Fatalf("Failed to serve plugin: %v", err)
	}
}
