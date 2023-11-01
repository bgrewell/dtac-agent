package plugins

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
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
		Ip:         "127.0.0.1",
		ApiVersion: "plug_api_1.0",
		encryptor:  utility.NewRpcEncryptor(key),
	}

	return plug, nil
}
