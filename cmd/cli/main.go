package main

import (
	"context"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/cli/commands"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

// TODO : --- DELETE BELOW ---

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the dtac-agent",
	Long:  `This command is used to set up the dtac-agent on the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up dtac-agent...")
		// Add your setup logic here
	},
}

//TODO : --- DELETE ABOVE ---

func loadConfig() (cfg *config.Configuration, err error) {
	logger, err := zap.NewProduction()
	return config.NewConfiguration(nil, logger)
}

func NewCommandLineInterface() *CommandLineInterface {

	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config: ", err)
		os.Exit(1)
	}

	cli := &CommandLineInterface{
		rootCmd: &cobra.Command{
			Use:   "dtac",
			Short: "dtac is a tool to configure the dtac-agent",
			Long:  `dtac is a command-line application tool to configure the dtac-agent on systems.`,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				ctx := context.WithValue(cmd.Context(), "config", cfg)
				cmd.SetContext(ctx)
			},
			Run: func(cmd *cobra.Command, args []string) {
				// Do nothing here
			},
		},
		config: cfg,
	}

	// Setup config commands
	commands.ConfigCmd.AddCommand(commands.ConfigViewCmd)
	commands.ConfigCmd.AddCommand(commands.ConfigEditCmd)

	// Setup setup commands

	// Setup root commands
	cli.rootCmd.AddCommand(commands.ConfigCmd)

	return cli
}

type CommandLineInterface struct {
	rootCmd        *cobra.Command
	configFilename string
	config         *config.Configuration
}

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
