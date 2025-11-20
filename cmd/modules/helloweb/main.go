package main

import (
	"github.com/bgrewell/dtac-agent/cmd/modules/helloweb/hellowebmodule"
	"log"

	"github.com/bgrewell/dtac-agent/pkg/modules"
)

func main() {

	m := hellowebmodule.NewHelloWebModule()

	h, err := modules.NewModuleHost(m)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
