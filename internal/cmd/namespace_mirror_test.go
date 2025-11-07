package cmd_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/cmd"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/anselstetter/git-build-number/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestNamespaceMirror(t *testing.T) {
	silence := bytes.NewBuffer([]byte{})

	t.Run("without flags confirm", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		reader := bufio.NewReader(strings.NewReader("y"))

		c := cmd.NewNamespaceMirrorCommand(bn, logger, reader)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()

		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
		assert.True(t, len(stdout.String()) > 0, "should ask for confirmation")
		assert.Equal(t, "", stderr.String())
	})
	t.Run("without flags no confirmation", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		reader := bufio.NewReader(strings.NewReader("n"))

		c := cmd.NewNamespaceMirrorCommand(bn, logger, reader)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()

		assert.NoError(t, err)
		assert.True(t, len(stdout.String()) > 0, "should ask for confirmation")
		assert.Equal(t, "", stderr.String())
	})
	t.Run("--yes", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		reader := bufio.NewReader(strings.NewReader("irrelevant"))

		c := cmd.NewNamespaceMirrorCommand(bn, logger, reader)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"--yes"})

		err := c.Execute()

		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
