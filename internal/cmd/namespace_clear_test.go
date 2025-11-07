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

func TestNamespaceClear(t *testing.T) {
	silence := bytes.NewBuffer([]byte{})

	t.Run("without flags confirm", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)
		bn.Set("test", "user", "email@domain.tld", 1)  // nolint:errcheck
		bn.Set("other", "user", "email@domain.tld", 4) // nolint:errcheck

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		reader := bufio.NewReader(strings.NewReader("y"))

		c := cmd.NewNamespaceClearCommand(bn, logger, reader)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()
		assert.NoError(t, err)

		for _, namespace := range []string{"test", "other"} {
			ns, err := bn.Get(namespace, "", "", false)
			assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
			assert.Nil(t, ns)
		}

		assert.True(t, len(stdout.String()) > 0, "should ask for confirmation")
		assert.Equal(t, "", stderr.String())
	})
	t.Run("without flags no confirmation", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)
		bn.Set("test", "user", "email@domain.tld", 1)  // nolint:errcheck
		bn.Set("other", "user", "email@domain.tld", 4) // nolint:errcheck

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		reader := bufio.NewReader(strings.NewReader("n"))

		c := cmd.NewNamespaceClearCommand(bn, logger, reader)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()
		assert.NoError(t, err)

		for _, namespace := range []string{"test", "other"} {
			ns, err := bn.Get(namespace, "", "", false)
			assert.NoError(t, err)
			assert.NotNil(t, ns)
		}

		assert.True(t, len(stdout.String()) > 0, "should ask for confirmation")
		assert.Equal(t, "", stderr.String())
	})
	t.Run("--yes", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)
		bn.Set("test", "user", "email@domain.tld", 1)  // nolint:errcheck
		bn.Set("other", "user", "email@domain.tld", 4) // nolint:errcheck

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		reader := bufio.NewReader(strings.NewReader("y"))

		c := cmd.NewNamespaceClearCommand(bn, logger, reader)
		c.SetOut(silence)
		c.SetErr(silence)
		c.SetArgs([]string{"--yes"})

		err := c.Execute()
		assert.NoError(t, err)

		for _, namespace := range []string{"test", "other"} {
			ns, err := bn.Get(namespace, "", "", false)
			assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
			assert.Nil(t, ns)
		}

		assert.Equal(t, "", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}
