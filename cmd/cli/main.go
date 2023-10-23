package main

import (
	"fmt"
	"os"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/cli/commands"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func loadConfig() (cfg *config.Configuration, err error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return config.NewConfiguration(nil, logger)
}

// NewCommandLineInterface returns a new instance of the CommandLineInterface struct.
func NewCommandLineInterface() *CommandLineInterface {

	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config: ", err)
		os.Exit(1)
	}

	cli := &CommandLineInterface{
		rootCmd: commands.NewRootCmd(cfg),
		config:  cfg,
	}

	// Setup config commands
	cfgCmd := commands.NewConfigCmd()
	cfgViewCmd := commands.NewConfigViewCmd()
	cfgEditCmd := commands.NewConfigEditCmd()
	cfgCmd.AddCommand(cfgViewCmd)
	cfgCmd.AddCommand(cfgEditCmd)

	// Setup root commands
	cli.rootCmd.AddCommand(cfgCmd)
	cli.rootCmd.AddCommand(commands.NewTokenCmd())

	return cli
}

// CommandLineInterface is a struct that contains the root command and the configuration.
type CommandLineInterface struct {
	rootCmd *cobra.Command
	config  *config.Configuration
}

// Run executes the root command.
func (cli *CommandLineInterface) Run() error {
	return cli.rootCmd.Execute()
}

func main() {
	cli := NewCommandLineInterface()
	if err := cli.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
