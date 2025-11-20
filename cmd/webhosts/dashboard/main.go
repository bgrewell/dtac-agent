package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bgrewell/dtac-agent/pkg/webhost"
)

func main() {
	// Create the dashboard module
	dashboard := NewDashboardModule()

	// Create the host (automatically detects standalone vs managed mode)
	host := webhost.NewDefaultWebhostHost(dashboard)

	// Set up context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutdown signal received, stopping module...")
		cancel()
	}()

	// Run the host
	if err := host.Run(ctx); err != nil {
		log.Fatalf("Host failed: %v", err)
	}
}
