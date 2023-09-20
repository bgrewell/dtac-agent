package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/plugin/examples/hello/hello_plugin"
	"github.com/bgrewell/gin-plugins/host"
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
