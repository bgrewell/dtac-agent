package modules

import (
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
)

// Module interface that all modules must implement
type Module interface {
	Name() string
	Register(args *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error
	RootPath() string
	LoggingStream(stream api.ModuleService_LoggingStreamServer) error
}
