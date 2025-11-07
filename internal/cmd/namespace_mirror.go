package cmd

import (
	"io"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewNamespaceMirrorCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger, reader io.Reader) *cobra.Command {
	var (
		remote string
		yes    bool
	)
	cmd := &cobra.Command{
		Use:    "mirror",
		Short:  "Mirror all local namespaces",
		PreRun: IgnoreAdditonalArgs(logger.StderrWriter(), 1),
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			return MirrorNamespaces(buildNumber, logger, remote, yes, reader)
		}),
	}
	cmd.Flags().StringVarP(&remote, "remote", "r", "origin", "the remote")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "I know what I’m doing")

	return cmd
}

func MirrorNamespaces(buildNumber buildnumber.BuildNumber, logger logger.Logger, remote string, yes bool, reader io.Reader) error {
	if yes || confirm(logger, "All remote namespaces that are not present locally will be deleted!\nUse --yes to skip the confirmation prompt.\n", "Continue?", reader) {
		return buildNumber.Mirror(remote)
	}
	return nil
}
