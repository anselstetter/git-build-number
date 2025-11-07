package cmd_test

import (
	"bytes"
	"testing"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/cmd"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/anselstetter/git-build-number/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestNamespaceDelete(t *testing.T) {
	silence := bytes.NewBuffer([]byte{})

	t.Run("without flags", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewNamespaceDeleteCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()

		assert.ErrorIs(t, err, cmd.ErrMissingNamespace)
		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
	t.Run("with invalid namespace", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewNamespaceDeleteCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"invalid"})

		err := c.Execute()

		assert.ErrorIs(t, err, repository.ErrReferenceNotFound)
		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
	t.Run("with valid namespace", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)
		bn.Set("test", "user", "email@domain.tld", 1)  // nolint:errcheck
		bn.Set("other", "user", "email@domain.tld", 4) // nolint:errcheck

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewNamespaceDeleteCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"test"})

		err := c.Execute()
		assert.NoError(t, err)

		ns, err := bn.Get("test", "", "", false)
		assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
		assert.Nil(t, ns)

		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
