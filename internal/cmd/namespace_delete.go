package cmd

import (
	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewNamespaceDeleteCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <namespace>...",
		Short: "Delete a namespace",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrMissingNamespace
			}
			return nil
		},
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			return Delete(buildNumber, logger, args...)
		}),
	}
	return cmd
}

func Delete(buildNumber buildnumber.BuildNumber, logger logger.Logger, namespaces ...string) error {
	return buildNumber.Delete(namespaces...)
}
