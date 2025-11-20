package hellomodule

import (
	"encoding/json"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/modules"
	"reflect"
)

// This sets a non-existent variable to the interface type of module then attempts to assign
// a pointer to HelloModule to it. This isn't needed, but it's a good way to ensure that the
// HelloModule struct implements the Module interface. If there are missing functions, this
// will fail to compile.
var _ modules.Module = &HelloModule{}

// NewHelloModule is a constructor that returns a new instance of the HelloModule
func NewHelloModule() *HelloModule {
	// Create a new instance of the module
	hm := &HelloModule{
		ModuleBase: modules.ModuleBase{},
	}
	// Ensure we set our root path which will be appended to all routes
	hm.SetRootPath("hello")

	// Return the new instance
	return hm
}

// HelloModule is the module struct that implements the Module interface
type HelloModule struct {
	// ModuleBase provides some helper functions
	modules.ModuleBase
}

// Name returns the name of the module type
// NOTE: this is intentionally not a pointer receiver otherwise it wouldn't work. This must be set at your module struct
// level. otherwise it will return the type of the ModuleBase struct instead.
func (h HelloModule) Name() string {
	t := reflect.TypeOf(h)
	return t.Name()
}

// Register registers the module with the module manager
func (h *HelloModule) Register(request *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error {
	*reply = api.ModuleRegisterResponse{
		ModuleType:   "basic",
		Capabilities: []string{"logging"},
	}

	// Convert the config json to a map. If you have a specific configuration type you should unmarshal into that type
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return err
	}

	// Print out a log message
	h.Log(modules.LoggingLevelInfo, "hello module registered", map[string]string{"module_type": "basic"})

	// Return no error
	return nil
}
