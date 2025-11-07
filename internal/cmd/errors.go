package cmd

import "errors"

var (
	ErrMissingBuildNumber = errors.New("please provide a build number")
	ErrInvalidNumber      = errors.New("not a valid number")
	ErrMissingNamespace   = errors.New("please provide a namespace")
)
