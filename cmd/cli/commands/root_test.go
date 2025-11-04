package commands

import (
	"bytes"
	"testing"

	"github.com/bgrewell/dtac-agent/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func GetDefaultTestConfig() *config.Configuration {
	cfg := config.DefaultConfig()
	for k, v := range cfg {
		viper.SetDefault(k, v)
	}
	var c config.Configuration
	if err := viper.Unmarshal(&c); err != nil {
		return nil
	}
	return &c
}

func TestNewRootCmd(t *testing.T) {

	cfg := GetDefaultTestConfig()
	if cfg == nil {
		t.Errorf("failed to get default config")
	}
	cmd := NewRootCmd(cfg)

	// test that the command has the expected properties
	assert.Equal(t, "dtac", cmd.Use)
	assert.Equal(t, "dtac is a tool to configure the dtac-agent", cmd.Short)
	assert.Equal(t, "dtac is a command-line application tool to configure the dtac-agent on systems.", cmd.Long)

	// Execute the command and capture the output
	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{})
	cmd.Execute()

	// Ensure that the help displays in the output
	assert.NotEqual(t, "", actual.String())

}
