package main

import (
	"github.com/bgrewell/dtac-agent/cmd/plugins/domaincontroller/plugin"
	"log"

	"github.com/bgrewell/dtac-agent/pkg/plugins"
)

func main() {

	p := plugin.NewDomainControllerPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
