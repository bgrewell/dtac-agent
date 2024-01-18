package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/docker/plugin"
	"log"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
)

func main() {

	p := plugin.NewDockerPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
