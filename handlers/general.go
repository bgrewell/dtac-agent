package handlers

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/mods"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	reflectorPort = 9000
)

var (
	Routes     gin.RoutesInfo
	Info       BasicInfo
	Reflectors []mods.Reflector
	udpPingWorkers map[string]*mods.UdpPingWorker
	tcpPingWorkers map[string]*mods.TcpPingWorker
)

func init() {
	Info = BasicInfo{}
	Info.Update()

	Reflectors = make([]mods.Reflector, 0)
	udp := mods.UdpReflector{}
	udp.SetPort(reflectorPort)
	udp.Start()
	Reflectors = append(Reflectors, &udp)

	tcp := mods.TcpReflector{}
	tcp.SetPort(reflectorPort)
	tcp.Start()
	Reflectors = append(Reflectors, &tcp)

	udpPingWorkers = make(map[string]*mods.UdpPingWorker)
	tcpPingWorkers = make(map[string]*mods.TcpPingWorker)
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