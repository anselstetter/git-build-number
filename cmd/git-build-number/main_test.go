package main

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("usage", func(t *testing.T) {
		t.Parallel()

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		code := run([]string{}, stdout, stderr, buildInfo)

		assert.Equal(t, 0, code)
		assert.Equal(t, "", stderr.String())
		assert.Contains(t, stdout.String(), "Usage:")
	})
	t.Run("invalid command", func(t *testing.T) {
		t.Parallel()

		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		code := run([]string{"invalid"}, stdout, stderr, buildInfo)

		assert.Equal(t, 1, code)
		assert.Equal(t, "", stdout.String())
		assert.Contains(t, stderr.String(), "unknown command")
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
