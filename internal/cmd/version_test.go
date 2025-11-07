package cmd_test

import (
	"bytes"
	"runtime/debug"

	"testing"

	"github.com/anselstetter/git-build-number/internal/cmd"
	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/anselstetter/git-build-number/internal/version"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	silence := bytes.NewBuffer([]byte{})

	t.Run("without flags", func(t *testing.T) {
		t.Parallel()

		v := version.New(buildInfo, "Dev")

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		logger := logger.New(logger.WithStdout(stdout), logger.WithStderr(stderr))

		c := cmd.NewVersionCommand(v, logger)
		c.SetOut(silence)
		c.SetErr(silence)

		err := c.Execute()

		assert.NoError(t, err)
		assert.Equal(t, "Test\n", stdout.String())
		assert.Equal(t, "", stderr.String())
	})
}

func buildInfo() (info *debug.BuildInfo, ok bool) {
	buildInfo := &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "-ldflags", Value: "-s -w -X main.Version=Test -s"},
		},
	}
	return buildInfo, true
}
