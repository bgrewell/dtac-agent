package calculatorplugin

import (
	"encoding/json"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
	_ "net/http/pprof"
	"reflect"
	"strconv"
)

// CalculationRequest represents a request for calculation
type CalculationRequest struct {
	Operation string  `json:"operation"` // add, subtract, multiply, divide
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

// CalculationResponse represents the result of a calculation
type CalculationResponse struct {
	Result float64 `json:"result"`
	CalleePlugin string `json:"callee_plugin,omitempty"` // If result came from another plugin
}

// This ensures that CalculatorPlugin implements the Plugin interface
var _ plugins.Plugin = &CalculatorPlugin{}

// NewCalculatorPlugin creates a new instance of the CalculatorPlugin
func NewCalculatorPlugin() *CalculatorPlugin {
	cp := &CalculatorPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}
	cp.SetRootPath("calculator")
	return cp
}

// CalculatorPlugin demonstrates plugin-to-plugin communication
type CalculatorPlugin struct {
	plugins.PluginBase
}

// Name returns the name of the plugin
func (c CalculatorPlugin) Name() string {
	t := reflect.TypeOf(c)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (c *CalculatorPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Initialize the broker client for plugin-to-plugin communication
	if request.BrokerAddress != "" {
		err := c.InitializeBrokerClient(request.BrokerAddress)
		if err != nil {
			c.Log(plugins.LevelWarning, "failed to initialize broker client", map[string]string{
				"error": err.Error(),
			})
		} else {
			c.Log(plugins.LevelInfo, "broker client initialized", map[string]string{
				"address": request.BrokerAddress,
			})
		}
	}

	// Declare endpoints
	authz := endpoint.AuthGroupAdmin.String()
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("calculate", endpoint.ActionCreate, "perform a calculation", c.Calculate, request.DefaultSecure, authz, endpoint.WithBody(&CalculationRequest{}), endpoint.WithOutput(&CalculationResponse{})),
		endpoint.NewEndpoint("calculate_via_hello", endpoint.ActionCreate, "demonstrate calling another plugin", c.CalculateViaHello, request.DefaultSecure, authz, endpoint.WithBody(&CalculationRequest{}), endpoint.WithOutput(&CalculationResponse{})),
		endpoint.NewEndpoint("list_plugins", endpoint.ActionRead, "list all available plugins", c.ListPlugins, request.DefaultSecure, authz),
	}

	// Register methods
	c.RegisterMethods(endpoints)

	// Convert to API endpoints
	for _, ep := range endpoints {
		aep := utility.ConvertEndpointToPluginEndpoint(ep)
		reply.Endpoints = append(reply.Endpoints, aep)
	}

	c.Log(plugins.LevelInfo, "calculator plugin registered", map[string]string{
		"endpoint_count": strconv.Itoa(len(endpoints)),
	})

	return nil
}

// Calculate performs a basic calculation
func (c *CalculatorPlugin) Calculate(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		var req CalculationRequest
		if err := json.Unmarshal(in.Body, &req); err != nil {
			return nil, fmt.Errorf("invalid request body: %w", err)
		}

		var result float64
		switch req.Operation {
		case "add":
			result = req.A + req.B
		case "subtract":
			result = req.A - req.B
		case "multiply":
			result = req.A * req.B
		case "divide":
			if req.B == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			result = req.A / req.B
		default:
			return nil, fmt.Errorf("unknown operation: %s", req.Operation)
		}

		response := CalculationResponse{
			Result: result,
		}

		return json.Marshal(response)
	}, "calculation result")
}

// CalculateViaHello demonstrates calling another plugin (hello) before performing calculation
// This shows how plugins can interact with each other
func (c *CalculatorPlugin) CalculateViaHello(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		var req CalculationRequest
		if err := json.Unmarshal(in.Body, &req); err != nil {
			return nil, fmt.Errorf("invalid request body: %w", err)
		}

		// First, call the hello plugin to demonstrate plugin-to-plugin communication
		c.Log(plugins.LevelInfo, "calling hello plugin before calculation", map[string]string{
			"operation": req.Operation,
		})

		helloResp, err := c.CallPlugin("HelloPlugin", "hello", endpoint.ActionRead, &endpoint.Request{})
		if err != nil {
			c.Log(plugins.LevelWarning, "failed to call hello plugin", map[string]string{
				"error": err.Error(),
			})
			// Continue with calculation even if hello plugin call fails
		} else {
			c.Log(plugins.LevelInfo, "received response from hello plugin", map[string]string{
				"response_size": strconv.Itoa(len(helloResp.Value)),
			})
		}

		// Now perform the actual calculation
		var result float64
		switch req.Operation {
		case "add":
			result = req.A + req.B
		case "subtract":
			result = req.A - req.B
		case "multiply":
			result = req.A * req.B
		case "divide":
			if req.B == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			result = req.A / req.B
		default:
			return nil, fmt.Errorf("unknown operation: %s", req.Operation)
		}

		response := CalculationResponse{
			Result:       result,
			CalleePlugin: "HelloPlugin",
		}

		return json.Marshal(response)
	}, "calculation result via hello plugin")
}

// ListPlugins demonstrates listing all available plugins
func (c *CalculatorPlugin) ListPlugins(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		pluginList, err := c.ListAvailablePlugins()
		if err != nil {
			return nil, fmt.Errorf("failed to list plugins: %w", err)
		}

		return json.Marshal(map[string]interface{}{
			"plugins": pluginList,
			"count":   len(pluginList),
		})
	}, "available plugins list")
}
