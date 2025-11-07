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

func TestFetch(t *testing.T) {
	silence := bytes.NewBuffer([]byte{})

	t.Run("without flags", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewFetchCommand(bn, logger)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()

		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
