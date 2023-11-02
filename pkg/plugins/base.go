package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"strings"
)

// PluginMethod declares the signature of plugin endpoint methods
type PluginMethod func(input *endpoint.InputArgs) (output *endpoint.ReturnVal, err error)

// PluginBase is a base struct that all plugins should embed as it implements the common shared methods
type PluginBase struct {
	Methods  map[string]endpoint.Func
	rootPath string
}

// Register is a default implementation of the Register method that must be implemented by the plugin therefor this one returns an error
func (p *PluginBase) Register(args RegisterArgs, reply *RegisterReply) error {
	return errors.New("this method must be implemented by the plugin")
}

// Call is a shim that calls the appropriate method on the plugin
func (p *PluginBase) Call(method string, args *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	key := strings.TrimPrefix(method, p.RootPath()+"/")
	if f, exists := p.Methods[key]; exists {
		return f(args)
	}

	return nil, fmt.Errorf("method %s not found", method)
}

// RegisterMethods is used to create the call map for the plugin
func (p *PluginBase) RegisterMethods(endpoints []*endpoint.Endpoint) {
	if p.Methods == nil {
		p.Methods = make(map[string]endpoint.Func)
	}
	for _, ep := range endpoints {
		p.Methods[ep.Path] = ep.Function
	}
}

// Name returns the name of the plugin
func (p *PluginBase) Name() string {
	return "UnnamedPlugin"
}

// RootPath returns the root path for the plugin
func (p *PluginBase) RootPath() string {
	return p.rootPath
}

// SetRootPath sets the value of rootPath for the plugin
func (p *PluginBase) SetRootPath(rootPath string) {
	p.rootPath = rootPath
}

// Serialize serializes the given interface to a string
func (p *PluginBase) Serialize(v interface{}) (string, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return "", e
	}
	return string(b), nil
}

// ToAPIEndpoint converts an endpoint to the API PluginEndpoint type
func ToAPIEndpoint(ep *endpoint.Endpoint) *PluginEndpoint {
	return &PluginEndpoint{
		Path:           ep.Path,
		Action:         ep.Action.String(),
		UsesAuth:       ep.UsesAuth,
		ExpectedArgs:   ep.ExpectedArgs,
		ExpectedBody:   ep.ExpectedBody,
		ExpectedOutput: ep.ExpectedOutput,
	}
}

// FromAPIEndpoint converts an API PluginEndpoint to an endpoint
func FromAPIEndpoint(ep *PluginEndpoint) *endpoint.Endpoint {
	action, _ := endpoint.ParseAction(ep.Action)
	return &endpoint.Endpoint{
		Path:           ep.Path,
		Action:         action,
		UsesAuth:       ep.UsesAuth,
		ExpectedArgs:   ep.ExpectedArgs,
		ExpectedBody:   ep.ExpectedBody,
		ExpectedOutput: ep.ExpectedOutput,
	}
}
