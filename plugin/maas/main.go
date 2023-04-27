package main

import (
	"github.com/BGrewell/dtac-agent/plugin/maas/maas_plugin"
	"github.com/bgrewell/gin-plugins/host"
	"log"
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
