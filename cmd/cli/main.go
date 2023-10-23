package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// TODO : --- DELETE BELOW ---
var rootCmd = &cobra.Command{
	Use:   "dtac",
	Short: "Dtac is a tool to configure the dtac-agent",
	Long:  `Dtac is a command-line application tool to configure the dtac-agent on systems.`,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Work with the config file",
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Viewing config from: ", viper.ConfigFileUsed())
	},
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Editing config key: ", args[0])
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup the dtac-agent",
	Long:  `This command is used to set up the dtac-agent on the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up dtac-agent...")
		// Add your setup logic here
	},
}

func init() {

	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(setupCmd)
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configEditCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
	}
	viper.ReadInConfig()
}

//TODO : --- DELETE ABOVE ---

func NewCommandLineInterface() *CommandLineInterface {
	cli := &CommandLineInterface{
		rootCmd: &cobra.Command{
			Use:   "dtac",
			Short: "Dtac is a CLI to configure the dtac-agent",
		},
	}
	cli.rootCmd.PersistentFlags().StringVar(&cli.cfgFilename, "config", "/etc/dtac/config.yaml", "config file (default is /etc/dtac/config.yaml)")

	return cli
}

type CommandLineInterface struct {
	rootCmd     *cobra.Command
	cfgFilename string
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
