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

// TestProxyRouteConfig_OAuthAuth tests OAuth 1.0a authentication (MAAS-style)
func TestProxyRouteConfig_OAuthAuth(t *testing.T) {
	expectedConsumerKey := "test-consumer-key"
	expectedToken := "test-token"
	expectedSecret := "test-secret"
	
	// Create a mock backend server that checks OAuth headers
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		
		// Check that it's an OAuth header
		if !strings.HasPrefix(auth, "OAuth ") {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "missing OAuth header",
			})
			return
		}
		
		// Verify OAuth components are present
		if !strings.Contains(auth, fmt.Sprintf("oauth_consumer_key=\"%s\"", expectedConsumerKey)) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid consumer key",
			})
			return
		}
		
		if !strings.Contains(auth, fmt.Sprintf("oauth_token=\"%s\"", expectedToken)) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid token",
			})
			return
		}
		
		// Check signature format (PLAINTEXT method: &{token_secret})
		expectedSig := fmt.Sprintf("%%26%s", expectedSecret)
		if !strings.Contains(auth, fmt.Sprintf("oauth_signature=\"%s\"", expectedSig)) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid signature",
				"auth":  auth,
			})
			return
		}
		
		// Verify required OAuth parameters
		requiredParams := []string{
			"oauth_signature_method=\"PLAINTEXT\"",
			"oauth_timestamp=",
			"oauth_nonce=",
			"oauth_version=\"1.0\"",
		}
		
		for _, param := range requiredParams {
			if !strings.Contains(auth, param) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": fmt.Sprintf("missing parameter: %s", param),
				})
				return
			}
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "authenticated",
		})
	}))
	defer backend.Close()

	// Create web module with OAuth proxy configuration
	webModule := &WebModuleBase{}
	webModule.SetRootPath("test")
	webModule.config = WebModuleConfig{
		Port:       0,
		StaticPath: "/",
		ProxyRoutes: []ProxyRouteConfig{
			{
				Name:      "maas",
				Target:    backend.URL,
				StripPath: true,
				AuthType:  "oauth",
				Credentials: ProxyCredentials{
					OAuthConsumerKey: expectedConsumerKey,
					OAuthToken:       expectedToken,
					OAuthTokenSecret: expectedSecret,
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
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/maas/machines/", webModule.GetPort()))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
		return
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

// TestRuntimeEnv_ConfigJS tests the /config.js endpoint
func TestRuntimeEnv_ConfigJS(t *testing.T) {
// Create web module with runtime environment variables
webModule := &WebModuleBase{}
webModule.SetRootPath("test")
webModule.config = WebModuleConfig{
Port:       0,
StaticPath: "/",
Debug:      true,
RuntimeEnv: map[string]string{
"API_URL":      "https://api.example.com",
"FEATURE_FLAG": "true",
"VERSION":      "1.0.0",
},
}

// Start the web module
err := webModule.Start()
if err != nil {
t.Fatalf("Failed to start web module: %v", err)
}
defer webModule.Stop()

// Give the server time to start
time.Sleep(100 * time.Millisecond)

// Test request to /config.js
resp, err := http.Get(fmt.Sprintf("http://localhost:%d/config.js", webModule.GetPort()))
if err != nil {
t.Fatalf("Failed to make request: %v", err)
}
defer resp.Body.Close()

// Check response status
if resp.StatusCode != http.StatusOK {
t.Errorf("Expected status 200, got %d", resp.StatusCode)
}

// Check content type
contentType := resp.Header.Get("Content-Type")
if !strings.Contains(contentType, "application/javascript") {
t.Errorf("Expected content type 'application/javascript', got '%s'", contentType)
}

// Check cache control headers
cacheControl := resp.Header.Get("Cache-Control")
if !strings.Contains(cacheControl, "no-cache") {
t.Errorf("Expected 'no-cache' in Cache-Control header, got '%s'", cacheControl)
}

// Read response body
body, err := io.ReadAll(resp.Body)
if err != nil {
t.Fatalf("Failed to read response body: %v", err)
}
bodyStr := string(body)

// Verify content contains the JavaScript object
if !strings.Contains(bodyStr, "window.__DTAC_CONFIG__") {
t.Errorf("Response should contain 'window.__DTAC_CONFIG__'")
}

// Verify each env var is present
if !strings.Contains(bodyStr, "API_URL") {
t.Errorf("Response should contain 'API_URL'")
}
if !strings.Contains(bodyStr, "https://api.example.com") {
t.Errorf("Response should contain 'https://api.example.com'")
}
if !strings.Contains(bodyStr, "FEATURE_FLAG") {
t.Errorf("Response should contain 'FEATURE_FLAG'")
}
if !strings.Contains(bodyStr, "VERSION") {
t.Errorf("Response should contain 'VERSION'")
}
}

// TestRuntimeEnv_ConfigJS_Empty tests /config.js with no runtime env
func TestRuntimeEnv_ConfigJS_Empty(t *testing.T) {
// Create web module without runtime environment variables
webModule := &WebModuleBase{}
webModule.SetRootPath("test")
webModule.config = WebModuleConfig{
Port:       0,
StaticPath: "/",
Debug:      true,
RuntimeEnv: nil, // No runtime env
}

// Start the web module
err := webModule.Start()
if err != nil {
t.Fatalf("Failed to start web module: %v", err)
}
defer webModule.Stop()

// Give the server time to start
time.Sleep(100 * time.Millisecond)

// Test request to /config.js
resp, err := http.Get(fmt.Sprintf("http://localhost:%d/config.js", webModule.GetPort()))
if err != nil {
t.Fatalf("Failed to make request: %v", err)
}
defer resp.Body.Close()

// Check response status
if resp.StatusCode != http.StatusOK {
t.Errorf("Expected status 200, got %d", resp.StatusCode)
}

// Read response body
body, err := io.ReadAll(resp.Body)
if err != nil {
t.Fatalf("Failed to read response body: %v", err)
}
bodyStr := string(body)

// Verify content contains empty object
if !strings.Contains(bodyStr, "window.__DTAC_CONFIG__ = {};") {
t.Errorf("Response should contain empty config object, got: %s", bodyStr)
}
}

// TestRuntimeEnv_ConfigJS_XSSPrevention tests that /config.js escapes dangerous values
func TestRuntimeEnv_ConfigJS_XSSPrevention(t *testing.T) {
// Create web module with potentially dangerous values
webModule := &WebModuleBase{}
webModule.SetRootPath("test")
webModule.config = WebModuleConfig{
Port:       0,
StaticPath: "/",
Debug:      true,
RuntimeEnv: map[string]string{
"DANGEROUS": "<script>alert(\"xss\")</script>",
"QUOTES":    "\"value with \"quotes\"",
"NEWLINES":  "line1\nline2",
},
}

// Start the web module
err := webModule.Start()
if err != nil {
t.Fatalf("Failed to start web module: %v", err)
}
defer webModule.Stop()

// Give the server time to start
time.Sleep(100 * time.Millisecond)

// Test request to /config.js
resp, err := http.Get(fmt.Sprintf("http://localhost:%d/config.js", webModule.GetPort()))
if err != nil {
t.Fatalf("Failed to make request: %v", err)
}
defer resp.Body.Close()

// Read response body
body, err := io.ReadAll(resp.Body)
if err != nil {
t.Fatalf("Failed to read response body: %v", err)
}
bodyStr := string(body)

// Verify dangerous characters are escaped
if strings.Contains(bodyStr, "<script>") {
t.Errorf("Response should not contain unescaped '<script>'")
}
if strings.Contains(bodyStr, "alert(\"xss\")") {
t.Errorf("Response should not contain unescaped alert")
}
// Check for proper escaping
if !strings.Contains(bodyStr, "\\u003c") {
t.Errorf("Response should contain escaped '<' character")
}
}

// TestParseWebModuleConfig_RuntimeEnv tests parsing of runtime_env configuration
func TestParseWebModuleConfig_RuntimeEnv(t *testing.T) {
configMap := map[string]interface{}{
"port":        8090,
"debug":       true,
"static_path": "/",
"runtime_env": map[string]interface{}{
"API_URL":      "https://api.example.com",
"FEATURE_FLAG": "true",
"NUMBER_VAL":   42, // Test non-string value conversion
},
}

config := ParseWebModuleConfig(configMap)

// Verify basic config
if config.Port != 8090 {
t.Errorf("Expected port 8090, got %d", config.Port)
}

// Verify runtime_env is parsed
if config.RuntimeEnv == nil {
t.Fatalf("RuntimeEnv should not be nil")
}

if config.RuntimeEnv["API_URL"] != "https://api.example.com" {
t.Errorf("Expected API_URL 'https://api.example.com', got '%s'", config.RuntimeEnv["API_URL"])
}

if config.RuntimeEnv["FEATURE_FLAG"] != "true" {
t.Errorf("Expected FEATURE_FLAG 'true', got '%s'", config.RuntimeEnv["FEATURE_FLAG"])
}

// Verify non-string value is converted to string
if config.RuntimeEnv["NUMBER_VAL"] != "42" {
t.Errorf("Expected NUMBER_VAL '42', got '%s'", config.RuntimeEnv["NUMBER_VAL"])
}
}

// TestInjectConfigScript tests the HTML injection function
func TestInjectConfigScript(t *testing.T) {
tests := []struct {
name     string
html     string
expected string
}{
{
name:     "inject before </head>",
html:     `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
expected: `<!DOCTYPE html><html><head><title>Test</title><script src="/config.js"></script>
</head><body></body></html>`,
},
{
name:     "inject after <head> if no </head>",
html:     `<!DOCTYPE html><html><head><title>Test</title><body></body></html>`,
expected: `<!DOCTYPE html><html><head>
<script src="/config.js"></script><title>Test</title><body></body></html>`,
},
{
name:     "inject after <body> if no head",
html:     `<!DOCTYPE html><html><body><p>content</p></body></html>`,
expected: `<!DOCTYPE html><html><body>
<script src="/config.js"></script><p>content</p></body></html>`,
},
{
name:     "case insensitive HEAD",
html:     `<!DOCTYPE html><html><HEAD><title>Test</title></HEAD><body></body></html>`,
expected: `<!DOCTYPE html><html><HEAD><title>Test</title><script src="/config.js"></script>
</HEAD><body></body></html>`,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
result := string(injectConfigScript([]byte(tt.html)))
if result != tt.expected {
t.Errorf("Expected:\n%s\n\nGot:\n%s", tt.expected, result)
}
})
}
}
