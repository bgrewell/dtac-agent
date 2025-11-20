package main

import (
	"github.com/bgrewell/dtac-agent/cmd/modules/hello/hellomodule"
	"log"

	"github.com/bgrewell/dtac-agent/pkg/modules"
)

func main() {

	m := hellomodule.NewHelloModule()

	h, err := modules.NewModuleHost(m)
	if err != nil {
		log.Fatal(err)
	}

	err = h.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
