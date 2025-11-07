package cmd

import (
	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewPushCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger) *cobra.Command {
	var (
		remote string
	)
	cmd := &cobra.Command{
		Use:    "push",
		Short:  "Push build number(s)",
		PreRun: IgnoreAdditonalArgs(logger.StderrWriter(), 1),
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			return Push(buildNumber, logger, remote)
		}),
	}
	cmd.Flags().StringVarP(&remote, "remote", "r", "origin", "the remote")

	return cmd
}

func Push(buildNumber buildnumber.BuildNumber, logger logger.Logger, remote string) error {
	err := buildNumber.Push(remote)
	if err != nil {
		return err
	}
	return err
}
