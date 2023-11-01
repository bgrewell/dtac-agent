package plugins

import (
	"context"
	"net/rpc"
)

// PluginInfo struct that contains information about a running plugin
type PluginInfo struct {
	Path         string
	Name         string
	RootPath     string
	Endpoints    []*PluginEndpoint
	Pid          int
	Proto        string
	Ip           string
	Port         int
	ApiVersion   string
	Key          []byte
	Rpc          *rpc.Client
	CancelToken  *context.CancelFunc
	ExitChan     chan int
	HasExited    bool
	ExitCode     int
	PluginConfig *PluginConfig
}
