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

func TestGet(t *testing.T) {
	silence := bytes.NewBuffer([]byte{})

	t.Run("without flags", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewGetCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()

		assert.ErrorIs(t, err, repository.ErrReferenceNotFound)
		assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
	t.Run("--create", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewGetCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"--create"})

		err := c.Execute()

		assert.NoError(t, err)
		assert.Equal(t, "1\n", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
	t.Run("--namespace", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)
		bn.Set("test", "user", "email@domain.tld", 123) // nolint:errcheck

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewGetCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"--namespace", "test"})

		err := c.Execute()

		assert.NoError(t, err)
		assert.Equal(t, "123\n", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
	t.Run("--user, --email", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)
		bn.Set("test", "user", "email@domain.tld", 123) // nolint:errcheck

		logger := logger.New(logger.WithStdout(silence), logger.WithStderr(silence))

		c := cmd.NewGetCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"--user", "First Last", "--email", "email@test.tld", "--create"})

		err := c.Execute()
		assert.NoError(t, err)

		commits, err := repo.Commits("refs/build-number/default")

		assert.NoError(t, err)
		assert.Equal(t, repository.Author{Name: "First Last", Email: "email@test.tld"}, commits[0].Author)
	})
}
