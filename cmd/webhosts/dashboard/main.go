package main

import (
	"log"

	"github.com/bgrewell/dtac-agent/pkg/webhost"
)

func main() {
	module := NewDashboardModule()
	host, err := webhost.NewHost(module)
	if err != nil {
		log.Fatalf("failed to create host: %v", err)
	}
	if err := host.Serve(); err != nil {
		log.Fatalf("host serve failed: %v", err)
	}
}
