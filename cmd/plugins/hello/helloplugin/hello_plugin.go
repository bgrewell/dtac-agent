package helloplugin

import (
	"encoding/json"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	_ "net/http/pprof" // Used for remote debugging of the plugin
	"reflect"
	"strconv"
)

// HelloMessage is just a simple helper struct to encapsulate the hello world message
type HelloMessage struct {
	Message string `json:"message"`
}

// This sets a non-existent variable to the interface type of plugin then attempts to assign
// a pointer to HelloPlugin to it. This isn't needed, but it's a good way to ensure that the
// HelloPlugin struct implements the Plugin interface. If there are missing functions, this
// will fail to compile.
var _ plugins.Plugin = &HelloPlugin{}

// NewHelloPlugin is a constructor that returns a new instance of the HelloPlugin
func NewHelloPlugin() *HelloPlugin {
	// Uncommenting the following anonymous function will allow remote debugging of the plugin by attaching a debugger
	// like the one built into goland. This is useful for debugging plugins.
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	// Create a new instance of the plugin
	hp := &HelloPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
		message: HelloMessage{
			Message: "this is an example of how to create a plugin. See the source at https://github.com/intel-innersource/frameworks.automation.dtac.agent/tree/main/plugin/examples/hello",
		},
	}
	// Ensure we set our root path which will be appended to all of our methods to help namespace them
	hp.SetRootPath("hello")

	// Return the new instance
	return hp
}

// HelloPlugin is the plugin struct that implements the Plugin interface
type HelloPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	message HelloMessage
}

// Name returns the name of the plugin type
// NOTE: this is intentionally not a pointer receiver otherwise it wouldn't work. This must be set at your plugin struct
// level. otherwise it will return the type of the PluginBase struct instead.
func (h HelloPlugin) Name() string {
	t := reflect.TypeOf(h)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (h *HelloPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Convert the config json to a map. If you have a specific configuration type you should unmarshal into that type
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return err
	}

	// Check if the configuration has the message set
	if message, ok := config["message"]; ok {
		h.message = HelloMessage{
			Message: message.(string),
		}
	}

	// Declare our endpoint(s)
	authz := endpoint.AuthGroupAdmin.String()
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("hello", endpoint.ActionRead, "this endpoint returns a hello world message", h.Hello, request.DefaultSecure, authz, endpoint.WithOutput(&HelloMessage{})),
	}

	// Register them with the plugin
	h.RegisterMethods(endpoints)

	// Convert to plugin endpoints and return
	for _, ep := range endpoints {
		aep := utility.ConvertEndpointToPluginEndpoint(ep)
		reply.Endpoints = append(reply.Endpoints, aep)
	}

	// Print out a log message
	h.Log(plugins.LevelInfo, "hello plugin registered", map[string]string{"endpoint_count": strconv.Itoa(len(endpoints))})

	// Return no error
	return nil
}

// Hello is the handler for the hello world route
func (h *HelloPlugin) Hello(in *endpoint.Request) (out *endpoint.Response, err error) {
	// Here we use the utility wrapper to help us add some additional context to the call and simplify the
	// code by having a helper function build the ReturnVal object for us.
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {h.Name()},
		}
		out, err := json.Marshal(h.message)
		if err != nil {
			return nil, nil, err
		}
		return headers, out, nil
	}, "hello plugin output message")
}
