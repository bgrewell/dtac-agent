package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/dsc/plugin"
	"log"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
)

func main() {

	p := plugin.NewDSCPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
