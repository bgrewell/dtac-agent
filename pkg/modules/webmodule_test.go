package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestProxyRouteConfig_NamedEndpoint tests that named endpoints work correctly
func TestProxyRouteConfig_NamedEndpoint(t *testing.T) {
	// Create a mock backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"path":   r.URL.Path,
			"query":  r.URL.RawQuery,
			"method": r.Method,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer backend.Close()

	// Create web module with proxy configuration
	webModule := &WebModuleBase{}
	webModule.SetRootPath("test")
	webModule.config = WebModuleConfig{
		Port:       0, // auto-assign
		StaticPath: "/",
		ProxyRoutes: []ProxyRouteConfig{
			{
				Name:      "maas",
				Target:    backend.URL,
				StripPath: true,
				AuthType:  "none",
			},
		},
		Debug: true,
	}

	// Start the web module
	err := webModule.Start()
	if err != nil {
		t.Fatalf("Failed to start web module: %v", err)
	}
	defer webModule.Stop()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Test request to the proxy endpoint using realistic MAAS API path
	// This simulates: GET /api/maas/machines/ which would proxy to the MAAS machines endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/maas/machines/?hostname=server1", webModule.GetPort()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the backend received the correct path (with trailing slash as MAAS expects)
	if result["path"] != "/machines/" {
		t.Errorf("Expected path '/machines/', got '%s'", result["path"])
	}

	// Verify query string was preserved
	if result["query"] != "hostname=server1" {
		t.Errorf("Expected query 'hostname=server1', got '%s'", result["query"])
	}
}

// TestProxyRouteConfig_BearerAuth tests bearer token authentication
func TestProxyRouteConfig_BearerAuth(t *testing.T) {
	expectedToken := "test-bearer-token"
	
	// Create a mock backend server that checks auth
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer "+expectedToken {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "authenticated",
		})
	}))
	defer backend.Close()

	// Create web module with bearer auth proxy configuration
	webModule := &WebModuleBase{}
	webModule.SetRootPath("test")
	webModule.config = WebModuleConfig{
		Port:       0,
		StaticPath: "/",
		ProxyRoutes: []ProxyRouteConfig{
			{
				Name:      "secure-api",
				Target:    backend.URL,
				StripPath: true,
				AuthType:  "bearer",
				Credentials: ProxyCredentials{
					Token: expectedToken,
				},
			},
		},
		Debug: true,
	}

	// Start the web module
	err := webModule.Start()
	if err != nil {
		t.Fatalf("Failed to start web module: %v", err)
	}
	defer webModule.Stop()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Test request to the proxy endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/secure-api/test", webModule.GetPort()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify authentication worked
	if result["status"] != "authenticated" {
		t.Errorf("Expected status 'authenticated', got '%s'", result["status"])
	}
}

// TestProxyRouteConfig_BasicAuth tests basic authentication
func TestProxyRouteConfig_BasicAuth(t *testing.T) {
	expectedUser := "testuser"
	expectedPass := "testpass"
	
	// Create a mock backend server that checks basic auth
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != expectedUser || pass != expectedPass {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "authenticated",
			"user":   user,
		})
	}))
	defer backend.Close()

	// Create web module with basic auth proxy configuration
	webModule := &WebModuleBase{}
	webModule.SetRootPath("test")
	webModule.config = WebModuleConfig{
		Port:       0,
		StaticPath: "/",
		ProxyRoutes: []ProxyRouteConfig{
			{
				Name:      "basic-api",
				Target:    backend.URL,
				StripPath: true,
				AuthType:  "basic",
				Credentials: ProxyCredentials{
					Username: expectedUser,
					Password: expectedPass,
				},
			},
		},
		Debug: true,
	}

	// Start the web module
	err := webModule.Start()
	if err != nil {
		t.Fatalf("Failed to start web module: %v", err)
	}
	defer webModule.Stop()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Test request to the proxy endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/basic-api/test", webModule.GetPort()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify authentication worked
	if result["status"] != "authenticated" {
		t.Errorf("Expected status 'authenticated', got '%s'", result["status"])
	}

	if result["user"] != expectedUser {
		t.Errorf("Expected user '%s', got '%s'", expectedUser, result["user"])
	}
}

// TestProxyRouteConfig_CustomHeaders tests custom header injection
func TestProxyRouteConfig_CustomHeaders(t *testing.T) {
	customHeaderKey := "X-Custom-Header"
	customHeaderValue := "custom-value"
	
	// Create a mock backend server that checks custom headers
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"custom_header": r.Header.Get(customHeaderKey),
		})
	}))
	defer backend.Close()

	// Create web module with custom headers
	webModule := &WebModuleBase{}
	webModule.SetRootPath("test")
	webModule.config = WebModuleConfig{
		Port:       0,
		StaticPath: "/",
		ProxyRoutes: []ProxyRouteConfig{
			{
				Name:      "custom-api",
				Target:    backend.URL,
				StripPath: true,
				AuthType:  "none",
				Credentials: ProxyCredentials{
					Headers: map[string]string{
						customHeaderKey: customHeaderValue,
					},
				},
			},
		},
		Debug: true,
	}

	// Start the web module
	err := webModule.Start()
	if err != nil {
		t.Fatalf("Failed to start web module: %v", err)
	}
	defer webModule.Stop()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Test request to the proxy endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/custom-api/test", webModule.GetPort()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify custom header was sent
	if result["custom_header"] != customHeaderValue {
		t.Errorf("Expected custom header '%s', got '%s'", customHeaderValue, result["custom_header"])
	}
}

// TestProxyRouteConfig_POSTRequest tests that POST requests work correctly
func TestProxyRouteConfig_POSTRequest(t *testing.T) {
	// Create a mock backend server that echoes POST data
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"echo":   string(body),
			"method": r.Method,
		})
	}))
	defer backend.Close()

	// Create web module with proxy configuration
	webModule := &WebModuleBase{}
	webModule.SetRootPath("test")
	webModule.config = WebModuleConfig{
		Port:       0,
		StaticPath: "/",
		ProxyRoutes: []ProxyRouteConfig{
			{
				Name:      "post-api",
				Target:    backend.URL,
				StripPath: true,
				AuthType:  "none",
			},
		},
		Debug: true,
	}

	// Start the web module
	err := webModule.Start()
	if err != nil {
		t.Fatalf("Failed to start web module: %v", err)
	}
	defer webModule.Stop()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Test POST request to the proxy endpoint
	testData := "test post data"
	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/post-api/test", webModule.GetPort()),
		"text/plain",
		strings.NewReader(testData),
	)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify echo
	if result["echo"] != testData {
		t.Errorf("Expected echo '%s', got '%s'", testData, result["echo"])
	}

	// Verify method
	if result["method"] != http.MethodPost {
		t.Errorf("Expected method 'POST', got '%s'", result["method"])
	}
}

// TestProxyRouteConfig_CustomPath tests using a custom path instead of default /api/<name>
func TestProxyRouteConfig_CustomPath(t *testing.T) {
	// Create a mock backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"path": r.URL.Path,
		})
	}))
	defer backend.Close()

	// Create web module with custom path
	webModule := &WebModuleBase{}
	webModule.SetRootPath("test")
	webModule.config = WebModuleConfig{
		Port:       0,
		StaticPath: "/static",
		ProxyRoutes: []ProxyRouteConfig{
			{
				Name:      "custom-path-api",
				Path:      "/custom/endpoint/",
				Target:    backend.URL,
				StripPath: true,
				AuthType:  "none",
			},
		},
		Debug: true,
	}

	// Start the web module
	err := webModule.Start()
	if err != nil {
		t.Fatalf("Failed to start web module: %v", err)
	}
	defer webModule.Stop()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Test request to the custom path
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/custom/endpoint/test", webModule.GetPort()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse response
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the backend received the correct path (stripped)
	if result["path"] != "/test" {
		t.Errorf("Expected path '/test', got '%s'", result["path"])
	}
}

// TestProxyRouteConfig_EdgeCases tests edge cases in path handling
func TestProxyRouteConfig_EdgeCases(t *testing.T) {
	// Create a mock backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"path": r.URL.Path,
		})
	}))
	defer backend.Close()

	// Test case 1: Empty path after stripping
	t.Run("EmptyPathAfterStripping", func(t *testing.T) {
		webModule := &WebModuleBase{}
		webModule.SetRootPath("test")
		webModule.config = WebModuleConfig{
			Port:       0,
			StaticPath: "/static",
			ProxyRoutes: []ProxyRouteConfig{
				{
					Name:      "edge-api",
					Path:      "/edge/",
					Target:    backend.URL,
					StripPath: true,
					AuthType:  "none",
				},
			},
			Debug: true,
		}

		err := webModule.Start()
		if err != nil {
			t.Fatalf("Failed to start web module: %v", err)
		}
		defer webModule.Stop()

		time.Sleep(100 * time.Millisecond)

		// Request to exactly the prefix - should result in / after stripping
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/edge/", webModule.GetPort()))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Should get "/" as the path
		if result["path"] != "/" {
			t.Errorf("Expected path '/', got '%s'", result["path"])
		}
	})

	// Test case 2: Path without leading slash after stripping
	t.Run("NoLeadingSlashAfterStripping", func(t *testing.T) {
		webModule := &WebModuleBase{}
		webModule.SetRootPath("test")
		webModule.config = WebModuleConfig{
			Port:       0,
			StaticPath: "/static",
			ProxyRoutes: []ProxyRouteConfig{
				{
					Name:      "slash-api",
					Path:      "/api/",
					Target:    backend.URL + "/base",
					StripPath: true,
					AuthType:  "none",
				},
			},
			Debug: true,
		}

		err := webModule.Start()
		if err != nil {
			t.Fatalf("Failed to start web module: %v", err)
		}
		defer webModule.Stop()

		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/users", webModule.GetPort()))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Should get "/base/users" as the path
		if result["path"] != "/base/users" {
			t.Errorf("Expected path '/base/users', got '%s'", result["path"])
		}
	})
}
