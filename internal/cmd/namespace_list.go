package cmd

import (
	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/spf13/cobra"
)

func NewNamespaceListCommand(buildNumber buildnumber.BuildNumber, logger logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "list",
		Short:  "List all namespaces",
		PreRun: IgnoreAdditonalArgs(logger.StderrWriter(), 1),
		RunE: SilenceUsageE(func(cmd *cobra.Command, args []string) error {
			return ListNamespaces(buildNumber, logger)
		}),
	}
	return cmd
}

func ListNamespaces(buildNumber buildnumber.BuildNumber, logger logger.Logger) error {
	namespaces, err := buildNumber.Namespaces()
	if err != nil {
		return err
	}
	out := []any{}
	for _, ns := range namespaces {
		out = append(out, ns.Name)
		out = append(out, ns.Entry.Number)
	}
	logger.StdoutTable(out...)
	return nil
}
