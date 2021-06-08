package handlers

import (
	"errors"
	"fmt"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/mods"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

func SendTimedUdpPingHandler(c *gin.Context) {
	start := time.Now()
	target := c.Param("target")
	if target == "" {
		WriteErrorResponseJSON(c, errors.New("missing target"))
		return
	}
	rtt, err := mods.UdpSendTimedPacket(target, reflectorPort, 2)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), rtt)
}

func SendTimedTcpPingHandler(c *gin.Context) {
	start := time.Now()
	target := c.Param("target")
	if target == "" {
		WriteErrorResponseJSON(c, errors.New("missing target"))
		return
	}
	rtt, err := mods.TcpSendTimedPacket(target, reflectorPort, 2)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, time.Since(start), rtt)
}

func GetUdpPingWorkerHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := udpPingWorkers[id]; ok {
		results := mods.PingOverview{
			Results: val.Results,
			Average: val.Average(),
			StdDev:  val.StdDev(),
		}
		WriteResponseJSON(c, time.Since(start), results)
		return
	}

	WriteErrorResponseJSON(c, fmt.Errorf("failed to find a udp ping worker with the id %s", id))
}

func GetTcpPingWorkerHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := tcpPingWorkers[id]; ok {
		results := mods.PingOverview{
			Results: val.Results,
			Average: val.Average(),
			StdDev:  val.StdDev(),
		}
		WriteResponseJSON(c, time.Since(start), results)
		return
	}

	WriteErrorResponseJSON(c, fmt.Errorf("failed to find a tcp ping worker with the id %s", id))
}

func DeleteUdpPingWorkerHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := udpPingWorkers[id]; ok {
		val.Stop()
		results := mods.PingOverview{
			Results: val.Results,
			Average: val.Average(),
			StdDev:  val.StdDev(),
		}
		WriteResponseJSON(c, time.Since(start), results)
		delete(udpPingWorkers, id)
		return
	}

	WriteErrorResponseJSON(c, fmt.Errorf("failed to find a udp ping worker with the id %s", id))
}

func DeleteTcpPingWorkerHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := tcpPingWorkers[id]; ok {
		val.Stop()
		results := mods.PingOverview{
			Results: val.Results,
			Average: val.Average(),
			StdDev:  val.StdDev(),
		}
		WriteResponseJSON(c, time.Since(start), results)
		delete(tcpPingWorkers, id)
		return
	}

	WriteErrorResponseJSON(c, fmt.Errorf("failed to find a tcp ping worker with the id %s", id))
}

func CreateUdpPingWorkerHandler(c *gin.Context) {
	start := time.Now()
	target := c.Param("target")
	var options *mods.UdpPingWorkerOptions
	if err := c.ShouldBindJSON(&options); err != nil {
		options = &mods.UdpPingWorkerOptions{}
	}
	options.Target = target
	if options.Port == 0 {
		options.Port = reflectorPort
	}
	if options.Interval == 0 {
		options.Interval = 30
	}
	if options.Timeout == 0 {
		options.Timeout = 2
	}

	id := uuid.New().String()
	log.WithFields(log.Fields{
		"target": target,
		"options": options,
		"id": id,
	}).Trace("creating new udp worker")
	w := mods.UdpPingWorker{}
	w.SetOptions(options)
	w.Start()
	udpPingWorkers[id] = &w
	WriteResponseJSON(c, time.Since(start), id)
}

func CreateTcpPingWorkerHandler(c *gin.Context) {
	start := time.Now()
	target := c.Param("target")
	var options *mods.TcpPingWorkerOptions
	if err := c.ShouldBindJSON(&options); err != nil {
		options = &mods.TcpPingWorkerOptions{}
	}
	options.Target = target
	if options.Port == 0 {
		options.Port = reflectorPort
	}
	if options.Interval == 0 {
		options.Interval = 30
	}
	if options.Timeout == 0 {
		options.Timeout = 2
	}

	id := uuid.New().String()
	log.WithFields(log.Fields{
		"target": target,
		"options": options,
		"id": id,
	}).Trace("creating new tcp worker")
	w := mods.TcpPingWorker{}
	w.SetOptions(options)
	w.Start()
	tcpPingWorkers[id] = &w
	WriteResponseJSON(c, time.Since(start), id)
}