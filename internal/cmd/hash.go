package cmd

import (
	"strconv"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewHashCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger) *cobra.Command {
	var (
		namespace string
	)
	cmd := &cobra.Command{
		Use:    "hash <number>",
		Short:  "Show the hash for a specific build number",
		PreRun: IgnoreAdditonalArgs(logger.StderrWriter(), 2),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrMissingBuildNumber
			}
			return nil
		},
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			number, _ := strconv.ParseInt(args[0], 10, 64)

			return Hash(buildNumber, logger, namespace, number)
		}),
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "the namespace")

	return cmd
}

func Hash(buildNumber buildnumber.BuildNumber, logger logger.Logger, namespace string, number int64) error {
	entry, err := buildNumber.Hash(namespace, number)
	if err != nil {
		return err
	}
	logger.Stdoutln(entry.Hash)
	return nil
}
