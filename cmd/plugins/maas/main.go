package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maas_plugin"
	"log"

	"github.com/bgrewell/gin-plugins/host"
)

func main() {

	p := new(maas_plugin.MAASPlugin)

	h, err := host.NewPluginHost(p, "this_is_not_a_security_feature")
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
