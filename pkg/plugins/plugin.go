package plugins

import (
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Register(args *api.RegisterRequest, reply *api.RegisterResponse) error
	Call(method string, args *endpoint.Request) (out *endpoint.Response, err error)
	RootPath() string
	LoggingStream(stream api.PluginService_LoggingStreamServer) error
	SetBroker(broker PluginBroker) // Allow the agent to inject the broker
	GetBroker() PluginBroker        // Allow plugins to access the broker
}
