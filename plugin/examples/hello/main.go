package main

import (
	"github.com/BGrewell/gin-plugins/host"
	"github.com/BGrewell/dtac-agent/plugin/examples/hello/hello_plugin"
	"log"
)

func main() {

	p := new(hello_plugin.HelloPlugin)

	h, err := host.NewPluginHost(p, "this_is_not_a_security_feature")
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
