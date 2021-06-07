package handlers

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/mods"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	Routes     gin.RoutesInfo
	Info       BasicInfo
	Reflectors []mods.Reflector
)

func init() {
	Info = BasicInfo{}
	Info.Update()

	Reflectors = make([]mods.Reflector, 0)
	udp := mods.UdpReflector{}
	udp.SetPort(9000)
	udp.Start()
	Reflectors = append(Reflectors, &udp)

	tcp := mods.TcpReflector{}
	tcp.SetPort(9000)
	tcp.Start()
	Reflectors = append(Reflectors, &tcp)
}

func SecretTestHandler(c *gin.Context) {
	user, err := AuthorizeUser(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
	}
	c.JSON(http.StatusOK, gin.H{"user": user.ID, "secret": "somesupersecretvalue"})
}

func HomeHandler(c *gin.Context) {
	// Update Routes
	start := time.Now()
	Info.UpdateRoutes(Routes)
	WriteResponseJSON(c, time.Since(start), Info)
}

func GetPingHandler(c *gin.Context) {
	c.Data(http.StatusOK, gin.MIMEPlain, []byte("pong"))
}

func GetReflectors(c *gin.Context) {
	start := time.Now()
	reflectors := make(map[string]int)
	for _, reflector := range Reflectors {
		reflectors[reflector.Proto()] = reflector.Port()
	}
	WriteResponseJSON(c, time.Since(start), reflectors)
}
