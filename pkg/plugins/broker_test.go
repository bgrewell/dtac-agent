package plugins

import (
	"testing"

	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"go.uber.org/zap"
)

// TestPluginBroker tests the basic functionality of the plugin broker
func TestPluginBroker(t *testing.T) {
	// Create a mock loader
	loader := &DefaultPluginLoader{
		plugins:  make(map[string]*PluginInfo),
		routeMap: make(map[string]*HandlerEntry),
		logger:   zap.NewNop(),
	}

	// Create the broker
	broker := NewPluginBroker(loader)

	// Test ListPlugins with empty plugin list
	plugins := broker.ListPlugins()
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}

	// Test IsPluginLoaded with non-existent plugin
	if broker.IsPluginLoaded("NonExistent") {
		t.Error("expected plugin to not be loaded")
	}
}

// TestBrokerClientCreation tests broker client creation
func TestBrokerClientCreation(t *testing.T) {
	// Test with empty address
	client, err := NewBrokerClient("")
	if err == nil {
		t.Error("expected error with empty broker address")
	}
	if client != nil {
		t.Error("expected nil client with empty broker address")
	}

	// Note: Testing actual connection would require a running gRPC server
	// which is beyond the scope of a unit test
}

// TestPluginBaseInitializeBrokerClient tests broker client initialization
func TestPluginBaseInitializeBrokerClient(t *testing.T) {
	pb := &PluginBase{}

	// Test with empty address (should not fail, just skip initialization)
	err := pb.InitializeBrokerClient("")
	if err != nil {
		t.Errorf("unexpected error with empty address: %v", err)
	}

	if pb.brokerClient != nil {
		t.Error("expected broker client to remain nil with empty address")
	}
}

// TestPluginBaseCallPluginWithoutBroker tests calling plugin without initialized broker
func TestPluginBaseCallPluginWithoutBroker(t *testing.T) {
	pb := &PluginBase{}

	// Try to call plugin without initializing broker
	_, err := pb.CallPlugin("TestPlugin", "method", endpoint.ActionRead, &endpoint.Request{})
	if err == nil {
		t.Error("expected error when calling plugin without broker client")
	}

	expectedMsg := "broker client is not initialized - call InitializeBrokerClient during plugin registration"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestPluginBaseListAvailablePluginsWithoutBroker tests listing plugins without initialized broker
func TestPluginBaseListAvailablePluginsWithoutBroker(t *testing.T) {
	pb := &PluginBase{}

	// Try to list plugins without initializing broker
	_, err := pb.ListAvailablePlugins()
	if err == nil {
		t.Error("expected error when listing plugins without broker client")
	}

	expectedMsg := "broker client is not initialized - call InitializeBrokerClient during plugin registration"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}
