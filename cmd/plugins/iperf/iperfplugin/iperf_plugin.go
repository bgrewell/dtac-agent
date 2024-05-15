package iperfplugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BGrewell/go-iperf"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"net/http"
	_ "net/http/pprof" // Used for remote debugging of the plugin
	"reflect"
	"strconv"
	"sync"
)

// This sets a non-existent variable to the interface type of plugin then attempts to assign
// a pointer to HelloPlugin to it. This isn't needed, but it's a good way to ensure that the
// HelloPlugin struct implements the Plugin interface. If there are missing functions, this
// will fail to compile.
var _ plugins.Plugin = &IperfPlugin{}

// NewIperfPlugin is a constructor that returns a new instance of the HelloPlugin
func NewIperfPlugin() *IperfPlugin {
	// Uncommenting the following anonymous function will allow remote debugging of the plugin by attaching a debugger
	// like the one built into goland. This is useful for debugging plugins. To do this do the following.
	// 1. Uncomment the anonymous function below
	// 2. Build the plugin like normal and execute via DTAC
	// 3.
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Create a new instance of the plugin
	p := &IperfPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
		iperfServerLock:  &sync.Mutex{},
		iperfClientLock:  &sync.Mutex{},
		iperfClients:     make(map[string]*iperf.Client),
		iperfServers:     make(map[string]*iperf.Server),
		iperfLiveResults: make(map[string]<-chan *iperf.StreamIntervalReport),
	}

	// Return the new instance
	return p
}

// IperfPlugin is the plugin struct that implements the Plugin interface
type IperfPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	bindAddr         string
	iperfClients     map[string]*iperf.Client
	iperfServers     map[string]*iperf.Server
	iperfServerLock  *sync.Mutex
	iperfClientLock  *sync.Mutex
	iperfLiveResults map[string]<-chan *iperf.StreamIntervalReport
	iperfController  *iperf.Controller
}

// Name returns the name of the plugin type
// NOTE: this is intentionally not a pointer receiver otherwise it wouldn't work. This must be set at your plugin struct
// level. otherwise it will return the type of the PluginBase struct instead.
func (p IperfPlugin) Name() string {
	t := reflect.TypeOf(p)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (p *IperfPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Convert the config json to a map. If you have a specific configuration type you should unmarshal into that type
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return err
	}

	// Check if the configuration has a bind address
	if bind, ok := config["bind"]; ok {
		p.bindAddr, ok = bind.(string)
		if !ok {
			p.Log(plugins.LevelError, "failed to convert bind to string", map[string]string{"type": reflect.TypeOf(bind).String()})
		} else {
			p.Log(plugins.LevelInfo, "control bind addr set via configuration file", map[string]string{"bind": p.bindAddr})
		}
	} else {
		p.bindAddr = "0.0.0.0"
	}

	// Check for controller port and set up controller
	controllerPort := 8191
	if configPort, ok := config["control_port"]; ok {
		cPortFloat, ok := configPort.(float64)
		if !ok {
			p.Log(plugins.LevelError, "failed to convert control port to int", map[string]string{"type": reflect.TypeOf(configPort).String()})
		} else {
			controllerPort = int(cPortFloat)
			p.Log(plugins.LevelInfo, "control port set via configuration file", map[string]string{"port": strconv.Itoa(controllerPort)})
		}
	}

	c, err := iperf.NewController(controllerPort)
	if err != nil {
		p.Log(plugins.LevelError, "error creating new controller", map[string]string{"error": err.Error()})
		return err
	}
	p.iperfController = c

	// Declare our endpoint(s)
	authz := endpoint.AuthGroupOperator.String()
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("server/start", endpoint.ActionCreate, "this endpoint starts a iperf server", p.CreateIperfServer, request.DefaultSecure, authz,
			endpoint.WithParameters(&IperfServerStartRequest{}),
			endpoint.WithOutput(&iperf.Server{}),
		),
		endpoint.NewEndpoint("client/start", endpoint.ActionCreate, "this endpoint starts a iperf client", p.CreateIperfClient, request.DefaultSecure, authz,
			endpoint.WithParameters(&IperfClientStartRequest{}),
			endpoint.WithOutput(&iperf.Client{}),
		),
		endpoint.NewEndpoint("reset", endpoint.ActionDelete, "this endpoint resets the iperf server and client", p.ResetIperf, request.DefaultSecure, authz),
		endpoint.NewEndpoint("server/stop", endpoint.ActionDelete, "this endpoint stops a iperf server", p.StopIperfServer, request.DefaultSecure, authz,
			endpoint.WithParameters(&IperfTestIDRequest{}),
		),
		endpoint.NewEndpoint("client/stop", endpoint.ActionDelete, "this endpoint stops a iperf client", p.StopIperfClient, request.DefaultSecure, authz,
			endpoint.WithParameters(&IperfTestIDRequest{}),
		),
		endpoint.NewEndpoint("client/results", endpoint.ActionRead, "this endpoint gets the results of a iperf client", p.GetIperfClientResults, request.DefaultSecure, authz,
			endpoint.WithParameters(&IperfTestIDRequest{}),
		),
		//r.GET("/iperf/server/results/:id", handlers.GetIperfServerTestResultsHandler)
		// TODO: Not supported yet
		//endpoint.NewEndpoint("server/results", endpoint.ActionRead, "this endpoint gets the results of a iperf server", p.GetIperfServerResults, request.DefaultSecure, authz,
		//	endpoint.WithParameters(&IperfTestIDRequest{}),
		//),
		// TODO: Not supported yet because it is specific to the Gin REST implementation
		//r.GET("/iperf/client/live/:id", handlers.GetIperfClientTestLiveHandler)
		//endpoint.NewEndpoint("client/live", endpoint.ActionRead, "this endpoint gets the live results of a iperf client", p.GetIperfClientLive, request.DefaultSecure, authz,
		//	endpoint.WithParameters(&IperfTestIDRequest{}),
		//),
		//r.GET("/iperf/client/results/:id", handlers.GetIperfClientTestResultsHandler)
	}

	// Register them with the plugin
	p.RegisterMethods(endpoints)

	// Convert to plugin endpoints and return
	for _, ep := range endpoints {
		aep := utility.ConvertEndpointToPluginEndpoint(ep)
		reply.Endpoints = append(reply.Endpoints, aep)
	}

	// Print out a log message
	p.Log(plugins.LevelInfo, "iperf plugin registered", map[string]string{"endpoint_count": strconv.Itoa(len(endpoints))})

	// Return no error
	return nil
}

// CreateIperfServer creates a new iperf server instance
func (p *IperfPlugin) CreateIperfServer(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {

		s, err := p.iperfController.NewServer()
		if err != nil {
			return nil, err
		}

		err = s.Start()
		if err != nil {
			return nil, err
		}

		p.iperfServerLock.Lock()
		p.iperfServers[s.Id] = s
		p.iperfServerLock.Unlock()

		return json.Marshal(s)
	}, "creates a new iperf server instance")
}

// CreateIperfClient creates a new iperf client instance
func (p *IperfPlugin) CreateIperfClient(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		host := "127.0.0.1"
		if hostVals, ok := in.Parameters["host"]; ok {
			if len(hostVals) > 0 {
				host = hostVals[0]
			}
		} else {
			return nil, errors.New("'host' is a required parameter")
		}

		var options *iperf.ClientOptions
		if in.Body != nil {
			err := json.Unmarshal(in.Body, options)
			if err != nil {
				return nil, err
			}
		}

		cli, err := p.iperfController.NewClient(host)
		if err != nil {
			return nil, err
		}

		if options != nil {
			options.Port = cli.Options.Port
			cli.LoadOptions(options)
			cli.SetHost(host)
		}
		// TODO: We don't support live option right now otherwise we would set json to false and SetModeLive()
		err = cli.Start()
		if err != nil {
			return nil, err
		}

		p.iperfClientLock.Lock()
		p.iperfClients[cli.Id] = cli
		p.iperfClientLock.Unlock()

		return json.Marshal(cli)
	}, "creates a new iperf client instance")
}

// ResetIperf resets all iperf client and server test instances
func (p *IperfPlugin) ResetIperf(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		servers := 0
		for key, value := range p.iperfServers {
			value.Stop()
			p.iperfServerLock.Lock()
			delete(p.iperfServers, key)
			p.iperfServerLock.Unlock()
			servers++
		}

		clients := 0
		for key, value := range p.iperfClients {
			value.Stop()
			p.iperfClientLock.Lock()
			delete(p.iperfClients, key)
			p.iperfClientLock.Unlock()
			clients++
		}
		return json.Marshal(map[string]int{
			"servers": servers,
			"clients": clients,
		})
	}, "resets all iperf client and server test instances")
}

// StopIperfServer stops and iperf server test instance
func (p *IperfPlugin) StopIperfServer(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		if id, ok := in.Parameters["id"]; ok {
			if s, ok := p.iperfServers[id[0]]; ok {
				s.Stop()
				p.iperfServerLock.Lock()
				delete(p.iperfServers, id[0])
				p.iperfServerLock.Unlock()
				_ = p.iperfController.StopServer(id[0])
				return json.Marshal(s)
			}

			return nil, fmt.Errorf("the specified id %s was not found on the system", id)

		}

		return nil, errors.New("the parameter 'id' is required")

	}, "stops and iperf server test instance")
}

// StopIperfClient stops an iperf client test instance
func (p *IperfPlugin) StopIperfClient(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		if id, ok := in.Parameters["id"]; ok {
			if c, ok := p.iperfClients[id[0]]; ok {
				c.Stop()
				report := c.Report()
				p.iperfClientLock.Lock()
				delete(p.iperfClients, id[0])
				p.iperfClientLock.Unlock()
				return json.Marshal(report)
			}

			return nil, fmt.Errorf("the specified id %s was not found on the system", id)

		}

		return nil, errors.New("the parameter 'id' is required")

	}, "stops an iperf client test instance")
}

// GetIperfServerResults gets iperf server test results
func (p *IperfPlugin) GetIperfServerResults(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this method has not been implemented")
	}, "gets iperf server test results")
}

// GetIperfClientLive gets iperf client test results lives
func (p *IperfPlugin) GetIperfClientLive(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		return nil, errors.New("this method has not been implemented")
	}, "gets iperf client test results lives")
}

// GetIperfClientResults gets iperf client test results
func (p *IperfPlugin) GetIperfClientResults(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		if id, ok := in.Parameters["id"]; ok {
			if c, ok := p.iperfClients[id[0]]; ok {
				if c.Running {
					return nil, errors.New("report not ready, test is still running")
				}
				report := c.Report()
				return json.Marshal(report)
			}

			return nil, fmt.Errorf("the specified id %s was not found on the system", id)

		}

		return nil, errors.New("the parameter 'id' is required")

	}, "gets iperf client test results")
}
