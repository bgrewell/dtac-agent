package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maasplugin"
	"log"

	"github.com/bgrewell/gin-plugins/host"
)

func main() {

	p := new(maasplugin.MAASPlugin)

	h, err := host.NewPluginHost(p, "still_using_the_old_gin_plugins")
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
