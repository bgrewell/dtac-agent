package helloplugin

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
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
func (h *HelloPlugin) Register(args plugins.RegisterArgs, reply *plugins.RegisterReply) error {
	*reply = plugins.RegisterReply{Endpoints: make([]*plugins.PluginEndpoint, 0)}

	// Check if the configuration has the message set
	if message, ok := args.Config["message"]; ok {
		h.message = HelloMessage{
			Message: message.(string),
		}
	}

	// Declare our endpoint(s)
	endpoints := []*endpoint.Endpoint{
		{
			Path:           "hello",
			Action:         endpoint.ActionRead,
			UsesAuth:       args.DefaultSecure,
			ExpectedArgs:   nil,
			ExpectedBody:   nil,
			ExpectedOutput: &HelloMessage{},
			Function:       h.Hello,
		},
	}

	// Register them with the plugin
	h.RegisterMethods(endpoints)

	// Convert to plugin endpoints and return
	for _, ep := range endpoints {
		reply.Endpoints = append(reply.Endpoints, plugins.ToAPIEndpoint(ep))
	}

	// Print out a log message
	h.Log(plugins.LevelInfo, "hello plugin registered", map[string]string{"endpoint_count": strconv.Itoa(len(endpoints))})

	// Return no error
	return nil
}

// Hello is the handler for the hello world route
func (h *HelloPlugin) Hello(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	// Here we use the utility wrapper to help us add some additional context to the call and simplify the
	// code by having a helper function build the ReturnVal object for us.
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, interface{}, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {h.Name()},
		}

		return headers, h.message, nil
	}, "hello plugin output message")
}
