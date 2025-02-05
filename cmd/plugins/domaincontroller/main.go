package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/domaincontroller/plugin"
	"log"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
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
