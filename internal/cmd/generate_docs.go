package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const synopsisGenerateDocs = `Generate docs for all commands in one of the following formats:

Markdown (md)
Man (man)
Yaml (yaml)`

func NewGenerateDocsCommand() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:    "generate-docs <dir>",
		Short:  "Generate docs",
		Long:   synopsisGenerateDocs,
		Hidden: true,
		Args: func(cmd *cobra.Command, args []string) error {
			validFormats := []string{"md", "man", "yaml"}

			if !slices.Contains(validFormats, format) {
				return fmt.Errorf("please use one of these formats: %s", strings.Join(validFormats, ", "))
			}
			if len(args) == 0 {
				return fmt.Errorf("please provide a directory")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := args[0]
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
			cmd.Root().DisableAutoGenTag = true

			switch format {
			case "md":
				return doc.GenMarkdownTree(cmd.Root(), dir)
			case "man":
				header := &doc.GenManHeader{
					Title:   "MINE",
					Section: "3",
				}
				return doc.GenManTree(cmd.Root(), header, dir)
			case "yaml":
				return doc.GenYamlTree(cmd.Root(), dir)
			default:
				return fmt.Errorf("unknown format: %s", format)
			}
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "md", "the doc type (md, man, yaml)")

	return cmd
}
