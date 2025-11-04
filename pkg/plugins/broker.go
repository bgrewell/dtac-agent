package plugins

import (
	"fmt"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// PluginBroker provides an interface for plugins to communicate with each other
// This allows plugins to call methods on other loaded plugins without going through
// the external API layer
type PluginBroker interface {
	// Call invokes a method on another plugin
	// pluginName: the name of the target plugin
	// method: the method path (e.g., "hello" for the hello endpoint)
	// action: the action type (Read, Write, Create, Delete)
	// request: the request to send to the plugin
	Call(pluginName string, method string, action endpoint.Action, request *endpoint.Request) (*endpoint.Response, error)

	// ListPlugins returns a list of all currently loaded plugin names
	ListPlugins() []string

	// IsPluginLoaded checks if a plugin with the given name is currently loaded
	IsPluginLoaded(pluginName string) bool
}

// DefaultPluginBroker implements the PluginBroker interface
type DefaultPluginBroker struct {
	loader *DefaultPluginLoader
}

// NewPluginBroker creates a new plugin broker instance
func NewPluginBroker(loader *DefaultPluginLoader) PluginBroker {
	return &DefaultPluginBroker{
		loader: loader,
	}
}

// Call invokes a method on another plugin
func (b *DefaultPluginBroker) Call(pluginName string, method string, action endpoint.Action, request *endpoint.Request) (*endpoint.Response, error) {
	// Check if plugin exists and is loaded
	plugin, exists := b.loader.plugins[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' not found", pluginName)
	}

	if plugin.HasExited {
		return nil, fmt.Errorf("plugin '%s' has exited with code %d", pluginName, plugin.ExitCode)
	}

	// Build the method key as the plugin expects it
	methodKey := fmt.Sprintf("%s:%s", action, method)

	// Use the loader's CallShim infrastructure to make the call
	// We need to construct a dummy endpoint to use with CallShim
	dummyEndpoint := &endpoint.Endpoint{
		Path:   fmt.Sprintf("%s/%s", plugin.RootPath, method),
		Action: action,
	}

	// Use the existing CallShim method which handles all the gRPC communication
	response, err := b.loader.CallShim(dummyEndpoint, request)
	if err != nil {
		return nil, fmt.Errorf("failed to call plugin '%s' method '%s': %w", pluginName, methodKey, err)
	}

	return response, nil
}

// ListPlugins returns a list of all currently loaded plugin names
func (b *DefaultPluginBroker) ListPlugins() []string {
	names := make([]string, 0, len(b.loader.plugins))
	for name := range b.loader.plugins {
		names = append(names, name)
	}
	return names
}

// IsPluginLoaded checks if a plugin with the given name is currently loaded
func (b *DefaultPluginBroker) IsPluginLoaded(pluginName string) bool {
	plugin, exists := b.loader.plugins[pluginName]
	if !exists {
		return false
	}
	return !plugin.HasExited
}
