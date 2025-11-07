package logger

import (
	"fmt"
	"io"
	"strings"
)

type Logger struct {
	options options
}

func New(opts ...option) Logger {
	options := newOptions(opts...)

	return Logger{
		options: options,
	}
}

func (l Logger) StdoutWriter() io.Writer {
	return l.options.stdout
}

func (l Logger) StderrWriter() io.Writer {
	return l.options.stderr
}

func (l Logger) Stdoutf(format string, a ...any) {
	if _, err := fmt.Fprintf(l.options.stdout, format, a...); err != nil {
		_ = "ignore"
	}
}

func (l Logger) StdoutTable(a ...any) {
	if len(a)%2 != 0 {
		l.Stderrf("Call to StdoutTable missing final value\n")
		return
	}
	keys := []string{}
	values := []string{}
	maxLength := 0

	for i := 0; i < len(a); i += 2 {
		key := fmt.Sprintf("%+v", a[i])
		value := fmt.Sprintf("%+v", a[i+1])

		keys = append(keys, key)
		values = append(values, value)
	}
	for _, k := range keys {
		length := len(k)

		if length > maxLength {
			maxLength = length
		}
	}
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		value := values[i]
		padding := strings.Repeat(" ", maxLength-len(key))

		l.Stdoutf("%s%s %s\n", key, padding, value)
	}
}

func (l Logger) Stdoutln(a ...any) {
	if _, err := fmt.Fprintln(l.options.stdout, a...); err != nil {
		_ = "ignore"
	}
}

func (l Logger) Stderrf(format string, a ...any) {
	if _, err := fmt.Fprintf(l.options.stderr, format, a...); err != nil {
		_ = "ignore"
	}
}
