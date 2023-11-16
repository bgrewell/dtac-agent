package plugins

import (
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Register(args RegisterArgs, reply *RegisterReply) error
	Call(method string, args *endpoint.InputArgs) (out *endpoint.ReturnVal, err error)
	RootPath() string
	LoggingStream(stream api.PluginService_LoggingStreamServer) error
}
