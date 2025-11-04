package main

import (
	"github.com/bgrewell/dtac-agent/cmd/plugins/calculator/calculatorplugin"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"log"
)

func main() {
	p := calculatorplugin.NewCalculatorPlugin()

	h, err := plugins.NewPluginHost(p)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
