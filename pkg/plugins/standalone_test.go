package plugins

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// TestPlugin is a simple plugin for testing
type TestPlugin struct {
	PluginBase
}

func (t *TestPlugin) Name() string {
	return "TestPlugin"
}

func (t *TestPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("test", endpoint.ActionRead, "test endpoint", t.TestEndpoint, false, ""),
	}

	t.RegisterMethods(endpoints)
	return nil
}

func (t *TestPlugin) TestEndpoint(in *endpoint.Request) (out *endpoint.Response, err error) {
	data := map[string]string{"status": "ok"}
	body, _ := json.Marshal(data)
	return &endpoint.Response{
		Value:   body,
		Headers: map[string][]string{"Content-Type": {"application/json"}},
	}, nil
}

func TestStandaloneConfig(t *testing.T) {
	// Test default config
	config := NewStandaloneConfig()
	if config.Enabled {
		t.Error("Expected Enabled to be false by default")
	}
	if config.Protocol != "http" {
		t.Errorf("Expected Protocol to be 'http', got '%s'", config.Protocol)
	}
	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", config.Port)
	}
	if config.Host != "0.0.0.0" {
		t.Errorf("Expected Host to be '0.0.0.0', got '%s'", config.Host)
	}
}

func TestStandaloneOptions(t *testing.T) {
	// Test with options
	config := NewStandaloneConfig(
		WithStandalone(),
		WithProtocol("https"),
		WithPort(9090),
		WithHost("127.0.0.1"),
	)

	if !config.Enabled {
		t.Error("Expected Enabled to be true")
	}
	if config.Protocol != "https" {
		t.Errorf("Expected Protocol to be 'https', got '%s'", config.Protocol)
	}
	if config.Port != 9090 {
		t.Errorf("Expected Port to be 9090, got %d", config.Port)
	}
	if config.Host != "127.0.0.1" {
		t.Errorf("Expected Host to be '127.0.0.1', got '%s'", config.Host)
	}
}

func TestStandaloneTLSOptions(t *testing.T) {
	// Test TLS options
	config := NewStandaloneConfig(
		WithStandalone(),
		WithTLS("/path/to/cert", "/path/to/key"),
	)

	if config.Protocol != "https" {
		t.Errorf("Expected Protocol to be 'https' when TLS is set, got '%s'", config.Protocol)
	}
	if config.TLSCertPath != "/path/to/cert" {
		t.Errorf("Expected TLSCertPath to be '/path/to/cert', got '%s'", config.TLSCertPath)
	}
	if config.TLSKeyPath != "/path/to/key" {
		t.Errorf("Expected TLSKeyPath to be '/path/to/key', got '%s'", config.TLSKeyPath)
	}
}

func TestRESTPluginHost(t *testing.T) {
	// Create test plugin
	plugin := &TestPlugin{
		PluginBase: PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}
	plugin.SetRootPath("test")

	// Create config for a random high port to avoid conflicts
	config := NewStandaloneConfig(
		WithStandalone(),
		WithProtocol("http"),
		WithPort(18080), // Use high port to avoid conflicts
		WithHost("127.0.0.1"),
	)

	// Create REST host
	host, err := NewRESTPluginHost(plugin, config)
	if err != nil {
		t.Fatalf("Failed to create REST plugin host: %v", err)
	}

	// Start server in goroutine
	go func() {
		if err := host.Serve(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test health endpoint
	resp, err := http.Get("http://127.0.0.1:18080/health")
	if err != nil {
		t.Fatalf("Failed to get health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var healthData map[string]string
	if err := json.Unmarshal(body, &healthData); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	if healthData["status"] != "healthy" {
		t.Errorf("Expected status to be 'healthy', got '%s'", healthData["status"])
	}
	if healthData["plugin"] != "TestPlugin" {
		t.Errorf("Expected plugin to be 'TestPlugin', got '%s'", healthData["plugin"])
	}

	// Test plugin endpoint
	resp, err = http.Get("http://127.0.0.1:18080/test/test")
	if err != nil {
		t.Fatalf("Failed to get test endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	var testData map[string]string
	if err := json.Unmarshal(body, &testData); err != nil {
		t.Fatalf("Failed to parse test response: %v", err)
	}

	if testData["status"] != "ok" {
		t.Errorf("Expected status to be 'ok', got '%s'", testData["status"])
	}

	// Shutdown server
	if err := host.Shutdown(nil); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}
}

func TestNewPluginHostWithStandalone(t *testing.T) {
	// Create test plugin
	plugin := &TestPlugin{
		PluginBase: PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}

	// Test creating a standalone host
	host, err := NewPluginHost(plugin, WithStandalone(), WithPort(18081))
	if err != nil {
		t.Fatalf("Failed to create plugin host with standalone: %v", err)
	}

	// Verify it's a REST host
	if _, ok := host.(*RESTPluginHost); !ok {
		t.Error("Expected host to be RESTPluginHost")
	}
}

func TestNewPluginHostWithoutStandalone(t *testing.T) {
	// Create test plugin
	plugin := &TestPlugin{
		PluginBase: PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}

	// Test creating a default host (without standalone)
	host, err := NewPluginHost(plugin)
	if err != nil {
		t.Fatalf("Failed to create plugin host: %v", err)
	}

	// Verify it's a default host
	if _, ok := host.(*DefaultPluginHost); !ok {
		t.Error("Expected host to be DefaultPluginHost")
	}
}
