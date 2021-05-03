package handlers

import (
	"errors"
	"fmt"
	"github.com/BGrewell/go-iperf"
	. "github.com/BGrewell/system-api/common"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	iperfClients     map[string]*iperf.Client
	iperfServers     map[string]*iperf.Server
	iperfLiveResults map[string]<-chan *iperf.StreamIntervalReport
	iperfController  *iperf.Controller
)

func init() {
	var err error
	iperfClients = make(map[string]*iperf.Client)
	iperfServers = make(map[string]*iperf.Server)
	iperfLiveResults = make(map[string]<-chan *iperf.StreamIntervalReport)
	iperfController, err = iperf.NewController(8090) //TODO: Expose in configuration file
	if err != nil {
		log.Printf("[WARNING] unable to instantiate iperf controller: %v\n", err)
		log.Printf("[WARNING] iperf will be unavailable")
	}
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
		options = nil
	}
	cli, err := iperfController.NewClient(host)
	if err != nil {
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
		log.Fatalf("error starting: %v", err)
	}
	//todo: need to figure out how to auto-poll for results and store them after it is done so the user can get them via a GET call with the ID
	//todo: this also needs to support live streaming of the results
	iperfClients[cli.Id] = cli
	WriteResponseJSON(c, time.Since(start), cli)
}

func CreateIperfServerTestHandler(c *gin.Context) {
	start := time.Now()
	s, err := iperfController.NewServer()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	err = s.Start()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	iperfServers[s.Id] = s
	WriteResponseJSON(c, time.Since(start), s)
}

func DeleteIperfClientTestHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "this function has not been implemented", "time": time.Now().Format(time.RFC3339Nano)})
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
		delete(iperfServers, id)
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
		delete(iperfServers, key)
		servers++
	}

	clients := 0
	for key, value := range iperfClients {
		value.Stop()
		delete(iperfClients, key)
		clients++
	}

	WriteResponseJSON(c, time.Since(start), fmt.Sprintf("stopped %d servers and %d clients", servers, clients))
}