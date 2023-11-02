package plugins

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Register(args RegisterArgs, reply *RegisterReply) error
	Call(method string, args *endpoint.InputArgs) (out *endpoint.ReturnVal, err error)
	RootPath() string
}
