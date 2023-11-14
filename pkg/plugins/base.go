package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
	"strings"
)

// PluginMethod declares the signature of plugin endpoint methods
type PluginMethod func(input *endpoint.InputArgs) (output *endpoint.ReturnVal, err error)

// PluginBase is a base struct that all plugins should embed as it implements the common shared methods
type PluginBase struct {
	LogChan  chan LogMessage
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

// LoggingStream is a function that sets up the logging channel for plugins to use so that they can log messages back
// to the agent. In advanced cases this could be overridden by the plugin to implement its own handling of the logging
// stream but there likely isn't a good reason to do that.
func (p *PluginBase) LoggingStream(stream api.PluginService_LoggingStreamServer) error {
	if p.LogChan == nil {
		p.LogChan = make(chan LogMessage, 4096)
	}
	for {
		msg := <-p.LogChan
		fields := make([]*api.LogField, 0)
		for k, v := range msg.Fields {
			fields = append(fields, &api.LogField{
				Key:   k,
				Value: v,
			})
		}
		err := stream.Send(&api.LogMessage{
			Level:   api.LogLevel(msg.Level),
			Message: msg.Message,
			Fields:  fields,
		})
		if err != nil {
			return err
		}
	}
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

// Log logs a message to the logging channel
func (p *PluginBase) Log(level LoggingLevel, message string, fields map[string]string) {
	if p.LogChan == nil {
		p.LogChan = make(chan LogMessage, 4096)
	}
	p.LogChan <- LogMessage{
		Level:   level,
		Message: message,
		Fields:  fields,
	}
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
