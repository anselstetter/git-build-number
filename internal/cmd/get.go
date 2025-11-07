package cmd

import (
	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewGetCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger) *cobra.Command {
	var (
		namespace string
		user      string
		email     string
		create    bool
	)
	cmd := &cobra.Command{
		Use:    "get",
		Short:  "Get the latest build number",
		PreRun: IgnoreAdditonalArgs(logger.StderrWriter(), 1),
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			return Get(buildNumber, logger, namespace, user, email, create)
		}),
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "the namespace")
	cmd.Flags().StringVarP(&user, "user", "u", "build number", "the author name")
	cmd.Flags().StringVarP(&email, "email", "e", "not set", "the author email")
	cmd.Flags().BoolVarP(&create, "create", "c", false, "create if missing")

	return cmd
}

func Get(buildNumber buildnumber.BuildNumber, logger logger.Logger, namespace string, user string, email string, create bool) error {
	entry, err := buildNumber.Get(namespace, user, email, create)
	if err != nil {
		return err
	}
	logger.Stdoutln(entry.Number)
	return nil
}
