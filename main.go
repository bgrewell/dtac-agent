package main

import (
	"github.com/BGrewell/system-api/handlers"
	"github.com/BGrewell/system-api/httprouting"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Default Router
	r := gin.Default()

	// General Routes
	httprouting.AddGeneralHandlers(r)

	// OS Specific Routes
	httprouting.AddOSSpecificHandlers(r)

	// Before starting update the handlers Routes var
	handlers.Routes = r.Routes()

	log.Println("system-api server is running http://localhost:8080")
	if err := r.Run(); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
