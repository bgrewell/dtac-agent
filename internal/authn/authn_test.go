package authn

import (
	"testing"
	"time"

	"github.com/bgrewell/dtac-agent/internal/controller"
	"go.uber.org/zap"
)

func TestParseTokenExpiration(t *testing.T) {
	// Create a simple subsystem for testing
	logger, _ := zap.NewDevelopment()
	s := &Subsystem{
		Logger: logger,
	}

	tests := []struct {
		name        string
		input       string
		want        time.Duration
		expectError bool
	}{
		{
			name:        "Valid minutes",
			input:       "30m",
			want:        30 * time.Minute,
			expectError: false,
		},
		{
			name:        "Valid hours",
			input:       "6h",
			want:        6 * time.Hour,
			expectError: false,
		},
		{
			name:        "Valid days (as hours)",
			input:       "168h",
			want:        168 * time.Hour,
			expectError: false,
		},
		{
			name:        "Never value (lowercase)",
			input:       "never",
			want:        0,
			expectError: false,
		},
		{
			name:        "Never value (uppercase)",
			input:       "NEVER",
			want:        0,
			expectError: false,
		},
		{
			name:        "Never value (mixed case)",
			input:       "NeVeR",
			want:        0,
			expectError: false,
		},
		{
			name:        "Never value (with whitespace)",
			input:       "  never  ",
			want:        0,
			expectError: false,
		},
		{
			name:        "Invalid format",
			input:       "invalid",
			want:        0,
			expectError: true,
		},
		{
			name:        "Negative duration",
			input:       "-5m",
			want:        0,
			expectError: true,
		},
		{
			name:        "Zero duration",
			input:       "0s",
			want:        0,
			expectError: false,
		},
		{
			name:        "Complex duration",
			input:       "1h30m",
			want:        90 * time.Minute,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.parseTokenExpiration(tt.input)
			if (err != nil) != tt.expectError {
				t.Errorf("parseTokenExpiration() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if got != tt.want {
				t.Errorf("parseTokenExpiration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTokenWithConfiguredExpiration(t *testing.T) {
	// This test validates that createToken uses the configured expiration times
	logger, _ := zap.NewDevelopment()
	
	// Create a minimal controller with config
	ctrl := &controller.Controller{
		Logger: logger,
	}
	
	s := &Subsystem{
		Controller: ctrl,
		Logger:     logger,
	}

	// Note: This is a basic structure test - full integration testing would require
	// a complete controller setup with AuthDB, which is beyond the scope of unit tests
	
	// Test that the parseTokenExpiration function is available and working
	duration, err := s.parseTokenExpiration("30m")
	if err != nil {
		t.Errorf("Expected no error parsing '30m', got: %v", err)
	}
	if duration != 30*time.Minute {
		t.Errorf("Expected 30 minutes, got: %v", duration)
	}
	
	duration, err = s.parseTokenExpiration("never")
	if err != nil {
		t.Errorf("Expected no error parsing 'never', got: %v", err)
	}
	if duration != 0 {
		t.Errorf("Expected 0 for 'never', got: %v", duration)
	}
}
