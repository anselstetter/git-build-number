package cmd

import (
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/anselstetter/git-build-number/internal/version"
	"github.com/spf13/cobra"
)

func NewVersionCommand(version version.Version, logger logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  "Just prints the version and exits",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Stdoutln(version.Version())
		},
	}
	return cmd
}
