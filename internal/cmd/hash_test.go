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

func TestHash(t *testing.T) {
	silence := bytes.NewBuffer([]byte{})

	t.Run("without args", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewHashCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()

		assert.ErrorIs(t, err, cmd.ErrMissingBuildNumber)
		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
	t.Run("with valid build number", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)
		entry, _ := bn.Set("default", "user", "email@domain.tld", 123)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewHashCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"123"})

		err := c.Execute()

		assert.NoError(t, err)
		assert.Equal(t, entry.Hash+"\n", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
	t.Run("with invalid build number", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewHashCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"1"})

		err := c.Execute()

		assert.ErrorIs(t, err, repository.ErrReferenceNotFound)
		assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
