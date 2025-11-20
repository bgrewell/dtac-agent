package modules

import (
	"context"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
)

// ModuleInfo struct that contains information about a running module
type ModuleInfo struct {
	Path          string
	Name          string
	RootPath      string
	ModuleType    string
	Capabilities  []string
	Endpoints     []*api.PluginEndpoint
	Pid           int
	RPCProto      string
	Proto         string
	IP            string
	Port          int
	APIVersion    string
	ModuleOptions *Options
	RPC           api.ModuleServiceClient
	CancelToken   *context.CancelFunc
	ExitChan      chan int
	HasExited     bool
	ExitCode      int
	ModuleConfig  *ModuleConfig
}
