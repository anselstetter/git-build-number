package cmd

import (
	"bufio"
	"io"
	"strings"

	"github.com/anselstetter/git-build-number/internal/logger"
)

func confirm(logger logger.Logger, msg string, confirmation string, reader io.Reader) bool {
	buf := bufio.NewReader(reader)

	logger.Stdoutln(msg)
	for {
		logger.Stdoutf("%s (y/N) ", confirmation)
		text, _ := buf.ReadString('\n')
		answer := strings.ToLower(strings.TrimSuffix(text, "\n"))

		if answer == "y" {
			return true
		}
		if answer == "n" {
			return false
		}
	}
}
