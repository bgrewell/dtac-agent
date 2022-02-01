package main

import (
	"fmt"
	. "github.com/BGrewell/system-api/plugin"
	"github.com/gin-gonic/gin"
	"plugin"
)

func main() {

	library, err := plugin.Open("/home/ben/repos/system-api/cmd/plugin_test/plugins/hello.so")
	if err != nil {
		panic(err)
	}

	loader, err := library.Lookup("Load")
	if err != nil {
		panic(err)
	}

	plug := loader.(func() Plugin)()
	name := plug.Name()

	fmt.Printf("loading: %s\n", name)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"plugin": name,
		})
	})

	err = plug.Register(nil, r.Group(fmt.Sprintf("plugins/%s", name)))
	if err != nil {
		panic(err)
	}

	r.Run(":9090")

}
