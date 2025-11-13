package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestTokenExpirationConfigParsing(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"
	
	configContent := `
auth:
  admin: admin
  pass: testpass
  default_secure: true
  model: /etc/dtac/auth_model.conf
  policy: /etc/dtac/auth_policy.csv
  access_token_expiration: 30m
  refresh_token_expiration: never
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	
	// Setup viper to read from our test config
	viper.Reset()
	viper.SetConfigFile(configFile)
	
	// Set defaults
	for k, v := range DefaultConfig() {
		viper.SetDefault(k, v)
	}
	
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	
	// Unmarshal into Configuration
	var cfg Configuration
	err = viper.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}
	
	// Verify the values were loaded correctly
	if cfg.Auth.AccessTokenExpiration != "30m" {
		t.Errorf("Expected access_token_expiration to be '30m', got '%s'", cfg.Auth.AccessTokenExpiration)
	}
	
	if cfg.Auth.RefreshTokenExpiration != "never" {
		t.Errorf("Expected refresh_token_expiration to be 'never', got '%s'", cfg.Auth.RefreshTokenExpiration)
	}
	
	t.Logf("Successfully loaded config with access_token_expiration=%s, refresh_token_expiration=%s", 
		cfg.Auth.AccessTokenExpiration, cfg.Auth.RefreshTokenExpiration)
}

func TestDefaultTokenExpirationValues(t *testing.T) {
	// Create logger
	logger, _ := zap.NewDevelopment()
	
	// Create a minimal temporary config file without token expiration fields
	tmpDir := t.TempDir()
	configFile := tmpDir + "/config.yaml"
	
	configContent := `
auth:
  admin: admin
  pass: testpass
  default_secure: true
`
	
	err := os.WriteFile(configFile, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	
	// Setup viper
	viper.Reset()
	viper.SetConfigFile(configFile)
	
	// Set defaults
	for k, v := range DefaultConfig() {
		viper.SetDefault(k, v)
	}
	
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	
	// Unmarshal into Configuration
	v := viper.New()
	for k, val := range viper.AllSettings() {
		v.Set(k, val)
	}
	var cfg Configuration
	err = v.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}
	
	cfg.logger = logger
	
	// Verify defaults were applied
	expectedAccessExpiration := "15m"
	expectedRefreshExpiration := "168h"
	
	if cfg.Auth.AccessTokenExpiration != expectedAccessExpiration {
		t.Errorf("Expected default access_token_expiration to be '%s', got '%s'", 
			expectedAccessExpiration, cfg.Auth.AccessTokenExpiration)
	}
	
	if cfg.Auth.RefreshTokenExpiration != expectedRefreshExpiration {
		t.Errorf("Expected default refresh_token_expiration to be '%s', got '%s'", 
			expectedRefreshExpiration, cfg.Auth.RefreshTokenExpiration)
	}
	
	t.Logf("Successfully verified defaults: access_token_expiration=%s, refresh_token_expiration=%s", 
		cfg.Auth.AccessTokenExpiration, cfg.Auth.RefreshTokenExpiration)
}
