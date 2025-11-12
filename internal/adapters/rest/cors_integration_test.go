package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bgrewell/dtac-agent/internal/basic"
	"github.com/bgrewell/dtac-agent/internal/config"
	"github.com/bgrewell/dtac-agent/internal/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// mockController creates a minimal controller for testing
func mockControllerWithCORS(corsEnabled bool, allowedOrigins []string) *controller.Controller {
	logger, _ := zap.NewDevelopment()
	
	cfg := &config.Configuration{
		APIs: config.APIEntries{
			REST: config.RESTAPIEntry{
				Enabled: true,
				Port:    8180,
				CORS: config.CORSConfig{
					Enabled:          corsEnabled,
					AllowedOrigins:   allowedOrigins,
					AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
					ExposedHeaders:   []string{"Content-Length"},
					AllowCredentials: false,
					MaxAge:           3600,
				},
			},
		},
	}

	return &controller.Controller{
		Config: cfg,
		Logger: logger,
	}
}

func TestCORSMiddlewareIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		corsEnabled        bool
		allowedOrigins     []string
		requestOrigin      string
		expectCORSHeaders  bool
		expectAllowOrigin  string
		expectStatus       int
	}{
		{
			name:              "CORS disabled - no CORS headers",
			corsEnabled:       false,
			allowedOrigins:    []string{"*"},
			requestOrigin:     "https://frontend.example.com",
			expectCORSHeaders: false,
			expectStatus:      http.StatusOK,
		},
		{
			name:               "CORS enabled with wildcard - allows any origin",
			corsEnabled:        true,
			allowedOrigins:     []string{"*"},
			requestOrigin:      "https://frontend.example.com",
			expectCORSHeaders:  true,
			expectAllowOrigin:  "*",
			expectStatus:       http.StatusOK,
		},
		{
			name:               "CORS enabled with specific origin - allows matching origin",
			corsEnabled:        true,
			allowedOrigins:     []string{"https://frontend.example.com", "https://test.com"},
			requestOrigin:      "https://frontend.example.com",
			expectCORSHeaders:  true,
			expectAllowOrigin:  "https://frontend.example.com",
			expectStatus:       http.StatusOK,
		},
		{
			name:              "CORS enabled with specific origin - rejects non-matching origin",
			corsEnabled:       true,
			allowedOrigins:    []string{"https://frontend.example.com"},
			requestOrigin:     "https://malicious.com",
			expectCORSHeaders: false,
			expectStatus:      http.StatusForbidden, // CORS rejects with 403
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := mockControllerWithCORS(tt.corsEnabled, tt.allowedOrigins)
			tls := make(map[string]basic.TLSInfo)
			
			adapter, err := NewAdapter(ctrl, &tls)
			assert.NoError(t, err)
			assert.NotNil(t, adapter)

			restAdapter := adapter.(*Adapter)
			
			// Add a simple test endpoint
			restAdapter.router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})

			// Create a test request with Origin header
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Origin", tt.requestOrigin)
			req.Host = "api.example.com" // Set a different host than the origin
			w := httptest.NewRecorder()

			// Execute the request
			restAdapter.router.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectStatus, w.Code)

			// Verify CORS headers
			if tt.expectCORSHeaders {
				allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
				assert.NotEmpty(t, allowOrigin, "Expected CORS headers but Access-Control-Allow-Origin is empty")
				assert.Equal(t, tt.expectAllowOrigin, allowOrigin)
			} else if tt.expectStatus == http.StatusOK {
				// Only check for empty CORS headers if request succeeded
				allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
				assert.Empty(t, allowOrigin, "Did not expect CORS headers but Access-Control-Allow-Origin is present")
			}
		})
	}
}

func TestCORSPreflightRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := mockControllerWithCORS(true, []string{"https://frontend.example.com"})
	tls := make(map[string]basic.TLSInfo)
	
	adapter, err := NewAdapter(ctrl, &tls)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	restAdapter := adapter.(*Adapter)
	
	// Add a simple test endpoint
	restAdapter.router.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// Create a preflight OPTIONS request
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://frontend.example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	req.Host = "api.example.com" // Set a different host than the origin
	w := httptest.NewRecorder()

	// Execute the request
	restAdapter.router.ServeHTTP(w, req)

	// Preflight requests should return 204 No Content or 200 OK
	assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK,
		"Expected 204 or 200, got %d", w.Code)

	// Verify CORS preflight headers
	assert.Equal(t, "https://frontend.example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.NotEmpty(t, w.Header().Get("Access-Control-Max-Age"))
}
