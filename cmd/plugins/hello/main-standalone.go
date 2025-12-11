package main

import (
	"github.com/bgrewell/dtac-agent/cmd/plugins/hello/helloplugin"
	"log"

	"github.com/bgrewell/dtac-agent/pkg/plugins"
)

// This example demonstrates running a plugin in standalone REST mode.
// The plugin can be configured via environment variables:
//
// DTAC_STANDALONE=true              - Enable standalone mode
// DTAC_STANDALONE_PROTOCOL=http     - Protocol (http or https)
// DTAC_STANDALONE_PORT=8080         - Port to listen on
// DTAC_STANDALONE_HOST=0.0.0.0      - Host to bind to
// DTAC_STANDALONE_TLS_CERT=/path    - Path to TLS certificate (for HTTPS)
// DTAC_STANDALONE_TLS_KEY=/path     - Path to TLS key (for HTTPS)
//
// Example usage:
//   DTAC_STANDALONE=true DTAC_STANDALONE_PORT=8080 ./hello-standalone
//
// The plugin will be available at:
//   http://localhost:8080/health        - Health check
//   http://localhost:8080/hello/hello   - Plugin endpoint (GET)

func main() {
	p := helloplugin.NewHelloPlugin()

	// Option 1: Using environment variables only
	// Just enable standalone mode, rest will be read from ENV
	h, err := plugins.NewPluginHost(p, plugins.WithStandalone())
	if err != nil {
		log.Fatal(err)
	}

	// Option 2 (commented): Using explicit configuration with options pattern
	// h, err := plugins.NewPluginHost(p,
	// 	plugins.WithStandalone(),
	// 	plugins.WithProtocol("http"),
	// 	plugins.WithPort(8080),
	// 	plugins.WithHost("0.0.0.0"),
	// )

	// Option 3 (commented): HTTPS with TLS
	// h, err := plugins.NewPluginHost(p,
	// 	plugins.WithStandalone(),
	// 	plugins.WithPort(8443),
	// 	plugins.WithTLS("/path/to/cert.pem", "/path/to/key.pem"),
	// )

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
