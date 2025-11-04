package commands

import (
	"fmt"
	"os"

	"github.com/bgrewell/dtac-agent/cmd/cli/consts"
	"github.com/bgrewell/dtac-agent/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewConfigCmd returns a new instance of the config command for the dtac tool.
func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Work with the config file",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

// NewConfigViewCmd returns a new instance of the config view command for the dtac tool.
func NewConfigViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "View configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := cmd.Context().Value(consts.KeyConfig)
			if cfg == nil {
				cmd.ErrOrStderr().Write([]byte("Error: config not found"))
				return
			}

			var value = cfg
			if len(args) > 0 {
				var err error
				value, err = config.GetConfigValue(cfg.(*config.Configuration), args[0])
				if err != nil {
					cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
					return
				}
			}

			yamlData, err := yaml.Marshal(value)
			if err != nil {
				cmd.ErrOrStderr().Write([]byte("Error: " + err.Error()))
				return
			}

			cmd.OutOrStdout().Write(yamlData)
		},
	}
}

// NewConfigEditCmd returns a new instance of the config edit command for the dtac tool.
func NewConfigEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit configuration",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Editing config key: ", args[0])
			fmt.Println("This method is not yet implemented")
			os.Exit(1)
		},
	}
}
