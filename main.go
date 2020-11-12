package main

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/handlers"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type HomeResponse struct {
	Status string   `json:"system_status"`
	Time   string   `json:"request_time"`
	Routes []string `json:"routes"`
}

func HomeHandler(c *gin.Context) {
	h := &HomeResponse{
		Status: "OK",
		Time:   time.Now().Format(time.RFC3339Nano),
		Routes: []string{
			"/",
			"/network/interfaces",
			"/network/interfaces/names",
			"/network/interface/{name:str} or {idx:int}",
		},
	}
	WriteResponseJSON(c, h)
}

func main() {

	r := gin.Default()

	// GET Routes
	r.GET("/", HomeHandler)
	r.GET("/network/interfaces", handlers.GetInterfacesHandler)
	r.GET("/network/interfaces/names", handlers.GetInterfaceNamesHandler)
	r.GET("/network/interface/:name", handlers.GetInterfaceByNameHandler)

	// POST Routes
	r.POST("/login", handlers.LoginHandler)

	log.Println("system-api server is running http://localhost:8080")
	r.Run()

}
