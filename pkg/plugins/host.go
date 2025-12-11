package plugins

import (
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
)

// PluginHost is the interface for the plugin host
type PluginHost interface {
	Serve() error
	GetPort() int
}

// NewPluginHost creates a new PluginHost with optional standalone configuration
func NewPluginHost(plugin Plugin, opts ...StandaloneOption) (hostPlugin PluginHost, err error) {
	// Create standalone config with options
	standaloneConfig := NewStandaloneConfig(opts...)

	// If standalone mode is enabled, create a REST host
	if standaloneConfig.Enabled {
		return NewRESTPluginHost(plugin, standaloneConfig)
	}

	// Otherwise, create the default gRPC host
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
