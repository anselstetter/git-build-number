package cmd

import (
	"io"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewNamespaceClearCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger, reader io.Reader) *cobra.Command {
	var (
		yes bool
	)
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Delete all namespaces",
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			return Clear(buildNumber, logger, yes, reader)
		}),
	}
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "I know what I’m doing")

	return cmd
}

func Clear(buildNumber buildnumber.BuildNumber, logger logger.Logger, yes bool, reader io.Reader) error {
	if yes || confirm(logger, "All local namespaces will be deleted!\nUse --yes to skip the confirmation prompt.\n", "Continue?", reader) {
		return buildNumber.Clear()
	}
	return nil
}
