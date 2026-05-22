// Command osquery is the dtac plugin that wraps a private osqueryd
// subprocess. Run it in either of two modes:
//
//	# Embedded gRPC mode (default — the dtac-agent spawns it).
//	./osquery.plugin
//
//	# Standalone REST mode.
//	DTAC_STANDALONE=true DTAC_STANDALONE_PORT=8080 ./osquery.plugin
//
// Selection is driven entirely by the DTAC_STANDALONE environment variable
// (consumed inside pkg/plugins.NewStandaloneConfig). The plugin owns the
// lifecycle of its osqueryd child and tears it down on SIGINT/SIGTERM.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bgrewell/dtac-agent/cmd/plugins/osquery/osqueryplugin"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
)

func main() {
	p := osqueryplugin.NewOsqueryPlugin()

	host, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		p.Shutdown(ctx)
		os.Exit(0)
	}()

	if err := host.Serve(); err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		p.Shutdown(ctx)
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	p.Shutdown(ctx)
}
