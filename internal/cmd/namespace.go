package cmd

import (
	"github.com/spf13/cobra"
)

func NewNamespaceCommand(cmds ...*cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage namespaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(cmds...)

	return cmd
}
