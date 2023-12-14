package main

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/iperf/iperfplugin"
	"log"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
)

func main() {

	p := iperfplugin.NewIperfPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
