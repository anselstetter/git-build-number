package cmd

import (
	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewIncCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger) *cobra.Command {
	var (
		namespace string
		user      string
		email     string
		force     bool
	)
	cmd := &cobra.Command{
		Use:    "inc",
		Short:  "Increment the build number",
		PreRun: IgnoreAdditonalArgs(logger.StderrWriter(), 1),
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			return Inc(buildNumber, logger, namespace, user, email, force)
		}),
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "the namespace")
	cmd.Flags().StringVarP(&user, "user", "u", "build number", "the author name")
	cmd.Flags().StringVarP(&email, "email", "e", "not set", "the author email")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force")

	return cmd
}

func Inc(buildNumber buildnumber.BuildNumber, logger logger.Logger, namespace string, user string, email string, force bool) error {
	entry, updated, err := buildNumber.Inc(namespace, user, email, force)
	if err != nil {
		return err
	}
	if !updated {
		logger.Stderrf("build number already set\nuse --force to override\n")
	}
	logger.Stdoutln(entry.Number)
	return nil
}
