package commands

import (
	"context"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/cli/consts"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/spf13/cobra"
)

// NewRootCmd returns a new instance of the root command for the dtac tool.
// It takes a pointer to a Configuration struct as input and returns a pointer to a Cobra command.
// The root command is used to configure the dtac-agent on systems via a command-line interface.
func NewRootCmd(cfg *config.Configuration) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dtac",
		Short: "dtac is a tool to configure the dtac-agent",
		Long:  `dtac is a command-line application tool to configure the dtac-agent on systems.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := context.WithValue(cmd.Context(), consts.KeyConfig, cfg)
			cmd.SetContext(ctx)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	return rootCmd
}
