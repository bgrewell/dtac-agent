package commands

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/cli/consts"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
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
		cfg := cmd.Context().Value(consts.KeyConfig)
		if cfg == nil {
			fmt.Println("Error: config not found")
			return
		}

		var value = cfg
		if len(args) > 0 {
			var err error
			value, err = config.GetConfigValue(cfg.(*config.Configuration), args[0])
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}

		yamlData, err := yaml.Marshal(value)
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
		fmt.Println("This method is not yet implemented")
		os.Exit(1)
	},
}
