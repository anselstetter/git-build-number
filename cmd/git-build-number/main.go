package main

import (
	"io"
	"os"
	"runtime/debug"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/cmd"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/anselstetter/git-build-number/internal/repository"
	"github.com/anselstetter/git-build-number/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, debug.ReadBuildInfo))
}

func run(args []string, stdout io.Writer, stderr io.Writer, buildInfoFunc version.BuildInfoFunc) int {
	logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))
	repo, err := repository.NewGitRepository(".")
	if err != nil {
		return fail(logger, err, 1)
	}
	version := version.New(buildInfoFunc, "Dev")
	buildNumber := buildnumber.New(repo)
	root := cmd.NewRootCommand()

	root.AddCommand(
		cmd.NewVersionCommand(version, logger),
		cmd.NewGetCommand(buildNumber, logger),
		cmd.NewSetCommand(buildNumber, logger),
		cmd.NewIncCommand(buildNumber, logger),
		cmd.NewPushCommand(buildNumber, logger),
		cmd.NewFetchCommand(buildNumber, logger),
		cmd.NewHashCommand(buildNumber, logger),
		cmd.NewNamespaceCommand(
			cmd.NewNamespaceListCommand(buildNumber, logger),
			cmd.NewNamespaceDeleteCommand(buildNumber, logger),
			cmd.NewNamespaceMirrorCommand(buildNumber, logger, os.Stdin),
			cmd.NewNamespaceClearCommand(buildNumber, logger, os.Stdin),
		),
		cmd.NewGenerateDocsCommand(),
	)
	root.SetArgs(args)
	root.SetOut(stdout)
	root.SetErr(stderr)

	if err := root.Execute(); err != nil {
		return fail(logger, err, 1)
	}
	return 0
}

func fail(logger logger.Logger, err error, exitCode int) int {
	logger.Stderrf("%s\n", err.Error())
	return exitCode
}
