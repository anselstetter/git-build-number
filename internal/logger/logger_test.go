package logger_test

import (
	"bytes"
	"testing"

	"github.com/anselstetter/git-build-number/internal/logger"
	"github.com/stretchr/testify/assert"
)

func TestKVLog(t *testing.T) {
	t.Run("should format entries", func(t *testing.T) {
		t.Parallel()

		stdout := bytes.NewBuffer([]byte{})
		logger := logger.New(logger.WithStdout(stdout))

		s := struct {
			a string
			b int
			c float64
		}{
			a: "a",
			b: 123,
			c: 123.123,
		}

		logger.StdoutTable(
			float32(2.9), "First value",
			"Longest key in this slice", 123,
			"Third", s,
		)
		want := `2.9                       First value
Longest key in this slice 123
Third                     {a:a b:123 c:123.123}
`
		assert.Equal(t, stdout.String(), want)
	})

	t.Run("should return warning on uneven arguments on stderr", func(t *testing.T) {
		t.Parallel()

		stderr := bytes.NewBuffer([]byte{})
		logger := logger.New(logger.WithStderr(stderr))

		logger.StdoutTable("First", "Second", "Third")
		want := "Call to StdoutTable missing final value\n"

		assert.Equal(t, stderr.String(), want)
	})
}

func TestStdOutWriter(t *testing.T) {
	t.Parallel()

	stdout := bytes.NewBuffer([]byte{})
	logger := logger.New(logger.WithStdout(stdout))

	assert.Equal(t, logger.StdoutWriter(), stdout)
}

func TestStdErrWriter(t *testing.T) {
	t.Parallel()

	stderr := bytes.NewBuffer([]byte{})
	logger := logger.New(logger.WithStderr(stderr))

	assert.Equal(t, logger.StderrWriter(), stderr)
}

func TestStdOutln(t *testing.T) {
	t.Parallel()

	stdout := bytes.NewBuffer([]byte{})
	logger := logger.New(logger.WithStdout(stdout))

	logger.Stdoutln("test")

	assert.Equal(t, stdout.String(), "test\n")
}
