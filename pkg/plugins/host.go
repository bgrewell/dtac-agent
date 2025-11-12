package plugins

import (
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
)

// PluginHost is the interface for the plugin host
type PluginHost interface {
	Serve() error
	GetPort() int
}

// NewPluginHost creates a new PluginHost
func NewPluginHost(plugin Plugin) (hostPlugin PluginHost, err error) {

	key := utility.NewRandomSymmetricKey()
	plug := &DefaultPluginHost{
		Plugin:     plugin,
		Proto:      "tcp",
		IP:         "127.0.0.1",
		APIVersion: "plug_api_1.0",
		encryptor:  utility.NewRPCEncryptor(key),
	}

	return plug, nil
}
