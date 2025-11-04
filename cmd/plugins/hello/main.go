package main

import (
	"github.com/bgrewell/dtac-agent/cmd/plugins/hello/helloplugin"
	"log"

	"github.com/bgrewell/dtac-agent/pkg/plugins"
)

func main() {

	p := helloplugin.NewHelloPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
