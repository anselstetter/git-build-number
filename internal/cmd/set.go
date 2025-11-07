package cmd

import (
	"fmt"
	"strconv"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewSetCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger) *cobra.Command {
	var (
		namespace string
		user      string
		email     string
	)
	cmd := &cobra.Command{
		Use:    "set <number>",
		Short:  "Set the build number",
		PreRun: IgnoreAdditonalArgs(logger.StderrWriter(), 2),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return ErrMissingBuildNumber
			}
			_, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidNumber, args[0])
			}
			return nil
		},
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			number, _ := strconv.ParseInt(args[0], 10, 64)
			return Set(buildNumber, logger, namespace, user, email, number)
		}),
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "the namespace")
	cmd.Flags().StringVarP(&user, "user", "u", "build number", "the author name")
	cmd.Flags().StringVarP(&email, "email", "e", "not set", "the author email")

	return cmd
}

func Set(buildNumber buildnumber.BuildNumber, logger logger.Logger, namespace string, user string, email string, number int64) error {
	entry, err := buildNumber.Set(namespace, user, email, number)
	if err != nil {
		return err
	}
	logger.Stdoutln(entry.Number)
	return nil
}
