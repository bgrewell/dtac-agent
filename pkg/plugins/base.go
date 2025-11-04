package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// PluginMethod declares the signature of plugin endpoint methods
type PluginMethod func(input *endpoint.Request) (output *endpoint.Response, err error)

// PluginBase is a base struct that all plugins should embed as it implements the common shared methods
type PluginBase struct {
	LogChan      chan LogMessage
	Methods      map[string]endpoint.Func
	rootPath     string
	broker       PluginBroker // Legacy field for backward compatibility
	brokerClient *BrokerClient
}

// Register is a default implementation of the Register method that must be implemented by the plugin therefor this one returns an error
func (p *PluginBase) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	return errors.New("this method must be implemented by the plugin")
}

// Call is a shim that calls the appropriate method on the plugin
func (p *PluginBase) Call(method string, args *endpoint.Request) (out *endpoint.Response, err error) {
	key := method
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
		p.Methods[fmt.Sprintf("%s:%s", ep.Action, ep.Path)] = ep.Function
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

// SetBroker sets the plugin broker for this plugin
func (p *PluginBase) SetBroker(broker PluginBroker) {
	p.broker = broker
}

// GetBroker returns the plugin broker for this plugin
func (p *PluginBase) GetBroker() PluginBroker {
	return p.broker
}

// InitializeBrokerClient initializes the broker client with the given broker address
// This should be called during plugin registration
func (p *PluginBase) InitializeBrokerClient(brokerAddress string) error {
	if brokerAddress == "" {
		// Broker is not available, which is fine for plugins that don't need plugin-to-plugin communication
		return nil
	}

	client, err := NewBrokerClient(brokerAddress)
	if err != nil {
		return fmt.Errorf("failed to initialize broker client: %w", err)
	}

	p.brokerClient = client
	return nil
}

// CallPlugin is a helper method that plugins can use to call other plugins
// This is the recommended way for plugins to communicate with each other
func (p *PluginBase) CallPlugin(pluginName string, method string, action endpoint.Action, request *endpoint.Request) (*endpoint.Response, error) {
	if p.brokerClient == nil {
		return nil, fmt.Errorf("broker client is not initialized - call InitializeBrokerClient during plugin registration")
	}
	return p.brokerClient.CallPlugin(pluginName, method, action, request)
}

// ListAvailablePlugins returns a list of all plugins currently loaded in the agent
func (p *PluginBase) ListAvailablePlugins() ([]string, error) {
	if p.brokerClient == nil {
		return nil, fmt.Errorf("broker client is not initialized - call InitializeBrokerClient during plugin registration")
	}
	return p.brokerClient.ListPlugins()
}
