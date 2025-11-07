package logger

import (
	"io"
	"os"
)

type option func(opts *options)

type options struct {
	stdout io.Writer
	stderr io.Writer
}

func newOptions(option ...option) options {
	opts := options{
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
	for _, fn := range option {
		fn(&opts)
	}
	return opts
}

func WithStdout(stdout io.Writer) option {
	return func(opts *options) {
		opts.stdout = stdout
	}
}

func WithStderr(stderr io.Writer) option {
	return func(opts *options) {
		opts.stderr = stderr
	}
}
