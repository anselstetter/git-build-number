package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func IgnoreAdditonalArgs(stderr io.Writer, n int) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if len(args) < n {
			return
		}
		if _, err := fmt.Fprintf(stderr, "Ignoring additional args: %s\n\n", strings.Join(args[1:], ", ")); err != nil {
			_ = "ignore"
		}
	}
}

func SilenceUsageE(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return f(cmd, args)
	}
}
