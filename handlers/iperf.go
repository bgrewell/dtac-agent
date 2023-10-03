package handlers

import (
	"errors"
	"fmt"
	"github.com/BGrewell/go-iperf"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	. "github.com/intel-innersource/frameworks.automation.dtac.agent/common"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/configuration"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"sync"
	"time"
)

//TODO: Iperf should be a module not built-in. Refactor this before release

var (
	iperfClients     map[string]*iperf.Client
	iperfServers     map[string]*iperf.Server
	iperfServerLock  sync.Mutex
	iperfClientLock  sync.Mutex
	iperfLiveResults map[string]<-chan *iperf.StreamIntervalReport
	iperfController  *iperf.Controller
)

func init() {
	go func() {
		// Delayed initialization
		for configuration.Config == nil {
			time.Sleep(100 * time.Millisecond)
		}
		if configuration.Config.Subsystems.Iperf {
			var err error
			iperfServerLock = sync.Mutex{}
			iperfClientLock = sync.Mutex{}
			iperfClients = make(map[string]*iperf.Client)
			iperfServers = make(map[string]*iperf.Server)
			iperfLiveResults = make(map[string]<-chan *iperf.StreamIntervalReport)
			iperfController, err = iperf.NewController(8090) //TODO: Expose in configuration file
			if err != nil {
				log.Printf("[WARNING] unable to instantiate iperf controller: %v\n", err)
				log.Printf("[WARNING] iperf will be unavailable")
			}
		}
	}()
}

func GetIperfClientTestLiveHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}

	if val, ok := iperfLiveResults[id]; ok {
		cli := iperfClients[id]

		count := 0
		c.Stream(func(w io.Writer) bool {
			select {
			case report := <-val:
				c.Render(-1, sse.Event{
					Event: "iperf-interval",
					Id:    strconv.Itoa(count),
					Data:  report,
				})
				count++
				return true
			case <-time.After(time.Duration(cli.Interval())*time.Second + (100 * time.Millisecond)):
				c.SSEvent("timeout", "a timeout occured while trying to get report")
				return false
			}
		})

		return
	}

	WriteErrorResponseJSON(c, fmt.Errorf("live results not available for id %s", id))
}

func GetIperfClientTestResultsHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := iperfClients[id]; ok {
		if val.Running {
			WriteErrorResponseJSON(c, errors.New("report not ready, test is still running"))
			return
		}
		report := val.Report()
		WriteResponseJSON(c, time.Since(start), report)
		return
	}

	WriteErrorResponseJSON(c, fmt.Errorf("failed to find a client with the id %s", id))
}

func GetIperfServerTestResultsHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := iperfServers[id]; ok {
		if val.Running {
			WriteErrorResponseJSON(c, errors.New("report not ready, test is still running"))
			return
		}
		//report := val.Report()
		WriteResponseJSON(c, time.Since(start), gin.H{"report": "reporting not yet implemented on server side"})
		return
	}

	WriteErrorResponseJSON(c, fmt.Errorf("failed to find a server with the id %s", id))
}

func CreateIperfClientTestHandler(c *gin.Context) {
	start := time.Now()
	host := c.Param("host")
	fmt.Printf("host: %s\n", host)
	var options *iperf.ClientOptions
	if err := c.ShouldBindJSON(&options); err != nil {
		log.WithFields(log.Fields{
			"host":    host,
			"options": options,
			"err":     err,
		}).Error("error binding iperf client options")
		fmt.Printf("error binding iperf client options: %v\n", err)
		options = nil
	}
	cli, err := iperfController.NewClient(host)
	if err != nil {
		log.WithFields(log.Fields{
			"host":    host,
			"options": options,
			"err":     err,
		}).Error("error getting new iperf client")
		fmt.Printf("error getting new iperf client: %v\n", err)
		WriteErrorResponseJSON(c, err)
		return
	}
	if options != nil {
		options.Port = cli.Options.Port // override port with the server assigned port
		cli.LoadOptions(options)
		cli.SetHost(host)
	}
	if _, ok := c.GetQuery("live"); ok {
		cli.SetJSON(false)
		ch := cli.SetModeLive()
		iperfLiveResults[cli.Id] = ch
	}
	err = cli.Start()
	if err != nil {
		log.WithFields(log.Fields{
			"host":    host,
			"options": options,
			"err":     err,
		}).Error("error starting iperf client")
		fmt.Printf("error starting iperf client: %v\n", err)
		log.Fatalf("error starting: %v", err)
	}

	iperfClientLock.Lock()
	iperfClients[cli.Id] = cli
	iperfClientLock.Unlock()
	WriteResponseJSON(c, time.Since(start), cli)
}

func CreateIperfServerTestHandler(c *gin.Context) {
	start := time.Now()
	s, err := iperfController.NewServer()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error getting new iperf server")
		WriteErrorResponseJSON(c, err)
		return
	}
	err = s.Start()
	if err != nil {
		log.WithFields(log.Fields{
			"server": s,
			"err":    err,
		}).Error("error starting iperf server")
		WriteErrorResponseJSON(c, err)
		return
	}

	iperfServerLock.Lock()
	iperfServers[s.Id] = s
	iperfServerLock.Unlock()
	WriteResponseJSON(c, time.Since(start), s)
}

func DeleteIperfClientTestHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := iperfClients[id]; ok {
		val.Stop()
		report := val.Report()
		iperfClientLock.Lock()
		delete(iperfClients, id)
		iperfClientLock.Unlock()
		WriteResponseJSON(c, time.Since(start), report)
		return
	}

	WriteErrorResponseJSON(c, errors.New(fmt.Sprintf("the specified id %s was not found on the system", id)))
}

func DeleteIperfServerTestHandler(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	if id == "" {
		WriteErrorResponseJSON(c, errors.New("id is a required parameter"))
		return
	}
	if val, ok := iperfServers[id]; ok {
		val.Stop()
		iperfServerLock.Lock()
		delete(iperfServers, id)
		iperfServerLock.Unlock()
		iperfController.StopServer(id)
		WriteResponseJSON(c, time.Since(start), val)
		return
	}

	WriteErrorResponseJSON(c, errors.New(fmt.Sprintf("the specified id %s was not found on the system", id)))
}

func DeleteIperfResetHandler(c *gin.Context) {
	start := time.Now()
	servers := 0
	for key, value := range iperfServers {
		value.Stop()
		iperfServerLock.Lock()
		delete(iperfServers, key)
		iperfServerLock.Unlock()
		servers++
	}

	clients := 0
	for key, value := range iperfClients {
		value.Stop()
		iperfClientLock.Lock()
		delete(iperfClients, key)
		iperfClientLock.Unlock()
		clients++
	}

	WriteResponseJSON(c, time.Since(start), fmt.Sprintf("stopped %d servers and %d clients.", servers, clients))
}
