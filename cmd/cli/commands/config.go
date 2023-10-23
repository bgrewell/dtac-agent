package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Work with the config file",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var ConfigViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd.Context().Value("config")
		if cfg == nil {
			fmt.Println("Error: config not found")
			return
		}

		yamlData, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Println(string(yamlData))
	},
}

var ConfigEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Editing config key: ", args[0])
	},
}
