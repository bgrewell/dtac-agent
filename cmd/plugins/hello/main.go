package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/hello/helloplugin"
	"log"

	"github.com/bgrewell/gin-plugins/host"
)

func main() {

	p := new(helloplugin.HelloPlugin)

	h, err := host.NewPluginHost(p, "this_is_not_a_security_feature")
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
