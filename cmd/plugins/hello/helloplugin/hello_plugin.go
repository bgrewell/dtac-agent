package helloplugin

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"reflect"
	"strings"
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

// RootPath returns the root path for the plugin
func (h HelloPlugin) RootPath() string {
	return "hello"
}

// Name returns the name of the plugin
func (h HelloPlugin) Name() string {
	t := reflect.TypeOf(h)
	return t.Name()
}

// Call is a shim that strips the root path off of the method name and then uses the base Call function to perform the call
func (h *HelloPlugin) Call(method string, args *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	key := strings.TrimPrefix(method, h.RootPath()+"/")
	return h.PluginBase.Call(key, args)
}

// Register registers the plugin with the plugin manager
func (h *HelloPlugin) Register(args plugins.RegisterArgs, reply *plugins.RegisterReply) error {
	*reply = plugins.RegisterReply{Endpoints: make([]*plugins.PluginEndpoint, 0)}

	// Register our one hello world route
	h.message = HelloMessage{
		Message: "this is an example of how to create a plugin. See the source at https://github.com/intel-innersource/frameworks.automation.dtac.agent/tree/main/plugin/examples/hello",
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

	// Return no error
	return nil
}

// Hello is the handler for the hello world route
func (h *HelloPlugin) Hello(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	//return utility.PlugFuncWrapper(in, out, func() (interface{}, error) {
	//	return h.Serialize(h.message)
	//}, "hello plugin message")
	out = &endpoint.ReturnVal{}
	out.Value = h.message
	out.Headers = in.Headers
	out.Params = in.Params
	out.Headers["Alive"] = []string{"true"}
	return out, nil
}
