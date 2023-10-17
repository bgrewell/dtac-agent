package hello_plugin

import (
	"net/http"
	"reflect"

	plugins "github.com/bgrewell/gin-plugins"
)

// HelloMessage is just a simple helper struct to encapsulate the hello world message
type HelloMessage struct {
	Message string `json:"message"`
}

// Ensure that our type meets the requirements for being a plugin
var _ plugins.Plugin = &HelloPlugin{}

// HelloPlugin is the plugin struct that implements the Plugin interface
type HelloPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	message HelloMessage
}

// RouteRoot returns the root path for the plugin
func (h HelloPlugin) RouteRoot() string {
	return "hello"
}

// Name returns the name of the plugin
func (h HelloPlugin) Name() string {
	t := reflect.TypeOf(h)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (h *HelloPlugin) Register(args plugins.RegisterArgs, reply *plugins.RegisterReply) error {
	*reply = plugins.RegisterReply{Routes: make([]*plugins.Route, 1)}

	// Register our one hello world route
	h.message = HelloMessage{
		Message: "this is an example of how to create a plugin. See the source at https://github.com/intel-innersource/frameworks.automation.dtac.agent/tree/main/plugin/examples/hello",
	}

	r := &plugins.Route{
		Path:       "hello",
		Method:     http.MethodGet,
		HandleFunc: "Hello",
	}
	reply.Routes[0] = r

	// Return no error
	return nil
}

// Hello is the handler for the hello world route
func (h *HelloPlugin) Hello(args plugins.Args, c *string) error {
	v, e := h.Serialize(h.message)
	if e != nil {
		return e
	}
	*c = v
	return nil
}
