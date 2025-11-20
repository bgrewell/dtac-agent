package modules

import (
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// Module interface that all modules must implement
type Module interface {
	Name() string
	Register(args *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error
	Call(method string, args *endpoint.Request) (out *endpoint.Response, err error)
	RootPath() string
	LoggingStream(stream api.ModuleService_LoggingStreamServer) error
}
