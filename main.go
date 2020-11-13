package main

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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
		Routes: []string{ //todo: auto-populate
			"/",
			"/network/interfaces",
			"/network/interfaces/names",
			"/network/interface/{name:str} or {idx:int}",
		},
	}
	WriteResponseJSON(c, h)
}

func SecretTestHandler(c *gin.Context) {
	user, err := handlers.AuthorizeUser(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
	}
	c.JSON(http.StatusOK, gin.H{"user": user.ID, "secret": "somesupersecretvalue"})
}

func main() {

	r := gin.Default()

	// GET Routes
	r.GET("/", HomeHandler)
	r.GET("/network/interfaces", handlers.GetInterfacesHandler)
	r.GET("/network/interfaces/names", handlers.GetInterfaceNamesHandler)
	r.GET("/network/interface/:name", handlers.GetInterfaceByNameHandler)
	r.GET("/secret", SecretTestHandler)

	// POST Routes
	r.POST("/login", handlers.LoginHandler)

	log.Println("system-api server is running http://localhost:8080")
	r.Run()

}
