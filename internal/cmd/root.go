package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git-build-number",
		Short: "Manage build numbers within a Git repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
	}
	cmd.Root().CompletionOptions.DisableDefaultCmd = true
	return cmd
}
