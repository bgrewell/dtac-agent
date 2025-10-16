package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maasplugin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"log"
)

func main() {

	p := maasplugin.NewMAASPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
