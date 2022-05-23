package handlers

import (
	"errors"
	"fmt"
	"github.com/BGrewell/system-agent/mods"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	reflectorPort = 9000
)

var (
	Reflectors     []mods.Reflector
	udpPingWorkers map[string]*mods.UdpPingWorker
	tcpPingWorkers map[string]*mods.TcpPingWorker
)

func init() {
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

func GetPingHandler(c *gin.Context) {
	start := time.Now()
	type rs struct {
		Response string `json:"response"`
	}
	r := rs{
		Response: "pong",
	}
	WriteResponseJSON(c, time.Since(start), r)
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
	payload := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	rtt, err := mods.UdpSendTimedPacket(target, reflectorPort, 2, 32, &payload)
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
	payload := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	rtt, err := mods.TcpSendTimedPacket(target, reflectorPort, 2, 32, &payload)
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
	fmt.Printf("options: %v\n", options)
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
	if options.PayloadSize == 0 {
		options.PayloadSize = 10
	}

	id := uuid.New().String()
	log.WithFields(log.Fields{
		"target":  target,
		"options": options,
		"id":      id,
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
	if options.PayloadSize == 0 {
		options.PayloadSize = 10
	}

	id := uuid.New().String()
	log.WithFields(log.Fields{
		"target":  target,
		"options": options,
		"id":      id,
	}).Trace("creating new tcp worker")
	w := mods.TcpPingWorker{}
	w.SetOptions(options)
	w.Start()
	tcpPingWorkers[id] = &w
	WriteResponseJSON(c, time.Since(start), id)
}

func DeleteResetAllPingWorkersHandler(c *gin.Context) {
	start := time.Now()
	tcpWorkers := 0
	for key, value := range tcpPingWorkers {
		value.Stop()
		delete(tcpPingWorkers, key)
		tcpWorkers++
	}

	udpWorkers := 0
	for key, value := range udpPingWorkers {
		value.Stop()
		delete(udpPingWorkers, key)
		udpWorkers++
	}

	WriteResponseJSON(c, time.Since(start), fmt.Sprintf("stopped %d udp workers and %d tcp workers.", udpWorkers, tcpWorkers))
}
