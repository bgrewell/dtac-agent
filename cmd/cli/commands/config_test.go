package commands

import (
	"bytes"
	"context"
	"testing"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/cli/consts"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNewConfigCmd(t *testing.T) {
	cmd := NewConfigCmd()

	// test that the command has the expected properties
	assert.Equal(t, "config", cmd.Use)
	assert.Equal(t, "Work with the config file", cmd.Short)

	// test that the command has a Run function
	assert.NotNil(t, cmd.Run)

	// test that the command is a Cobra command
	assert.IsType(t, &cobra.Command{}, cmd)
}

func TestNewConfigViewCmd(t *testing.T) {
	cmd := NewConfigViewCmd()

	// Ensure we set the context manually
	cfg := GetDefaultTestConfig()

	ctx := context.WithValue(context.Background(), consts.KeyConfig, cfg)
	cmd.SetContext(ctx)

	if cmd.Use != "view" {
		t.Errorf("Expected cmd.Use to be 'view', but got '%s'", cmd.Use)
	}

	if cmd.Short != "View configuration" {
		t.Errorf("Expected cmd.Short to be 'View configuration', but got '%s'", cmd.Short)
	}

	// Test with no arguments
	actual := new(bytes.Buffer)
	actual.Reset()
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute expected no error, but got %v", err)
	}
	output, err := yaml.Marshal(cfg)
	if err != nil {
		t.Errorf("Marshal expected no error, but got %v", err)
	}

	expectedOutput := string(output)
	if actual.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, actual.String())
	}

	// Test with valid argument
	actual.Reset()
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"authn.user"})
	err = cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	expectedOutput = "admin\n"
	if actual.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, actual.String())
	}

	// Test with invalid argument
	actual.Reset()
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs([]string{"key.doesnt.exist"})
	err = cmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	expectedOutput = "Error: Key not found: key.doesnt.exist"
	if actual.String() != expectedOutput {
		t.Errorf("Expected output to be '%s', but got '%s'", expectedOutput, actual.String())
	}
}
