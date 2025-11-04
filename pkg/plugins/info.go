package plugins

import (
	"context"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
)

// PluginInfo struct that contains information about a running plugin
type PluginInfo struct {
	Path          string
	Name          string
	RootPath      string
	Endpoints     []*api.PluginEndpoint
	Pid           int
	RPCProto      string
	Proto         string
	IP            string
	Port          int
	APIVersion    string
	PluginOptions *Options
	RPC           api.PluginServiceClient
	CancelToken   *context.CancelFunc
	ExitChan      chan int
	HasExited     bool
	ExitCode      int
	PluginConfig  *PluginConfig
}
