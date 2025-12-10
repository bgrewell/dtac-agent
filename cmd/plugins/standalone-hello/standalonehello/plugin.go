package standalonehello

import (
	"encoding/json"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
	"reflect"
	"strconv"
)

// HelloMessage is a simple message structure
type HelloMessage struct {
	Message string `json:"message"`
	Plugin  string `json:"plugin"`
	Mode    string `json:"mode"`
}

// StandaloneHelloPlugin demonstrates a plugin that can run with its own REST interface
type StandaloneHelloPlugin struct {
	plugins.PluginBase
	message HelloMessage
}

// Ensure StandaloneHelloPlugin implements the Plugin interface
var _ plugins.Plugin = &StandaloneHelloPlugin{}

// NewStandaloneHelloPlugin creates a new instance of the standalone hello plugin
func NewStandaloneHelloPlugin() *StandaloneHelloPlugin {
	hp := &StandaloneHelloPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
		message: HelloMessage{
			Message: "Hello from standalone plugin! This plugin is running its own REST API without the DTAC agent.",
			Plugin:  "StandaloneHelloPlugin",
			Mode:    "standalone",
		},
	}
	hp.SetRootPath("hello")
	return hp
}

// Name returns the name of the plugin
func (h StandaloneHelloPlugin) Name() string {
	t := reflect.TypeOf(h)
	return t.Name()
}

// Register registers the plugin's endpoints
func (h *StandaloneHelloPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Parse config if provided
	if request.Config != "" && request.Config != "{}" {
		var config map[string]interface{}
		err := json.Unmarshal([]byte(request.Config), &config)
		if err == nil {
			if message, ok := config["message"]; ok {
				h.message.Message = message.(string)
			}
		}
	}

	// Define endpoints
	authz := endpoint.AuthGroupGuest.String()
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint(
			"message",
			endpoint.ActionRead,
			"Returns a hello message from the standalone plugin",
			h.GetMessage,
			request.DefaultSecure,
			authz,
			endpoint.WithOutput(&HelloMessage{}),
		),
		endpoint.NewEndpoint(
			"echo",
			endpoint.ActionCreate,
			"Echoes back the provided message",
			h.Echo,
			request.DefaultSecure,
			authz,
			endpoint.WithBody(&HelloMessage{}),
			endpoint.WithOutput(&HelloMessage{}),
		),
	}

	// Register methods
	h.RegisterMethods(endpoints)

	// Convert to API endpoints
	for _, ep := range endpoints {
		aep := utility.ConvertEndpointToPluginEndpoint(ep)
		reply.Endpoints = append(reply.Endpoints, aep)
	}

	// Log registration
	h.Log(plugins.LevelInfo, "standalone hello plugin registered", map[string]string{
		"endpoint_count": strconv.Itoa(len(endpoints)),
		"mode":           "standalone",
	})

	return nil
}

// GetMessage returns the hello message
func (h *StandaloneHelloPlugin) GetMessage(in *endpoint.Request) (*endpoint.Response, error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {h.Name()},
			"X-PLUGIN-MODE": {"standalone"},
		}
		out, err := json.Marshal(h.message)
		if err != nil {
			return nil, nil, err
		}
		return headers, out, nil
	}, "standalone hello plugin message")
}

// Echo echoes back the provided message
func (h *StandaloneHelloPlugin) Echo(in *endpoint.Request) (*endpoint.Response, error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		// Parse the input message
		var inputMsg HelloMessage
		if len(in.Body) > 0 {
			if err := json.Unmarshal(in.Body, &inputMsg); err != nil {
				return nil, nil, err
			}
		}

		// Create response with the echoed message
		response := HelloMessage{
			Message: inputMsg.Message,
			Plugin:  h.Name(),
			Mode:    "echo",
		}

		headers := map[string][]string{
			"X-PLUGIN-NAME": {h.Name()},
			"X-PLUGIN-MODE": {"standalone"},
		}

		out, err := json.Marshal(response)
		if err != nil {
			return nil, nil, err
		}
		return headers, out, nil
	}, "standalone hello plugin echo")
}
