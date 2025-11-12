package rest

import (
	"testing"

	"github.com/bgrewell/dtac-agent/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestCORSConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		corsConfig     config.CORSConfig
		expectEnabled  bool
		expectOrigins  []string
		expectMethods  []string
		expectHeaders  []string
	}{
		{
			name: "CORS disabled",
			corsConfig: config.CORSConfig{
				Enabled: false,
			},
			expectEnabled: false,
		},
		{
			name: "CORS enabled with wildcard",
			corsConfig: config.CORSConfig{
				Enabled:          true,
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST"},
				AllowedHeaders:   []string{"Content-Type"},
				AllowCredentials: false,
				MaxAge:           3600,
			},
			expectEnabled: true,
			expectOrigins: []string{"*"},
			expectMethods: []string{"GET", "POST"},
			expectHeaders: []string{"Content-Type"},
		},
		{
			name: "CORS enabled with specific origins",
			corsConfig: config.CORSConfig{
				Enabled:          true,
				AllowedOrigins:   []string{"https://example.com", "https://test.com"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
				AllowedHeaders:   []string{"Content-Type", "Authorization"},
				ExposedHeaders:   []string{"Content-Length"},
				AllowCredentials: true,
				MaxAge:           7200,
			},
			expectEnabled: true,
			expectOrigins: []string{"https://example.com", "https://test.com"},
			expectMethods: []string{"GET", "POST", "PUT", "DELETE"},
			expectHeaders: []string{"Content-Type", "Authorization"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectEnabled, tt.corsConfig.Enabled)
			if tt.expectEnabled {
				assert.Equal(t, tt.expectOrigins, tt.corsConfig.AllowedOrigins)
				assert.Equal(t, tt.expectMethods, tt.corsConfig.AllowedMethods)
				assert.Equal(t, tt.expectHeaders, tt.corsConfig.AllowedHeaders)
			}
		})
	}
}

func TestDefaultCORSConfig(t *testing.T) {
	// Test that default CORS configuration is sensible
	defaultCORS := config.CORSConfig{
		Enabled:          false,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           3600,
	}

	assert.False(t, defaultCORS.Enabled, "CORS should be disabled by default")
	assert.Contains(t, defaultCORS.AllowedOrigins, "*", "Default should allow all origins when enabled")
	assert.Contains(t, defaultCORS.AllowedMethods, "GET", "Default should include GET method")
	assert.Contains(t, defaultCORS.AllowedMethods, "POST", "Default should include POST method")
	assert.Contains(t, defaultCORS.AllowedHeaders, "Authorization", "Default should include Authorization header")
	assert.Equal(t, 3600, defaultCORS.MaxAge, "Default max age should be 3600 seconds")
}
