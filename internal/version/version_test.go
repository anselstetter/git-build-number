package version_test

import (
	"runtime/debug"
	"testing"

	"github.com/anselstetter/git-build-number/internal/version"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name              string
		module            debug.Module
		buildInfoSettings []debug.BuildSetting
		ok                bool
		expected          string
	}{
		{
			name:     "should display the default version",
			ok:       false,
			expected: "Dev",
		},
		{
			name:              "should display the version passed as ldflags",
			buildInfoSettings: []debug.BuildSetting{{Key: "-ldflags", Value: "-s -w -X main.Version=Test -other"}},
			ok:                true,
			expected:          "Test",
		},
		{
			name:              "should display the BuildInfo.Main.Version when version is missing from ldflags",
			module:            debug.Module{Version: "Module"},
			buildInfoSettings: []debug.BuildSetting{{Key: "-ldflags", Value: "-s -w -X -other"}},
			ok:                true,
			expected:          "Module",
		},
		{
			name:     "should display the BuildInfo.Main.Version when ldflags are missing",
			module:   debug.Module{Version: "Module"},
			ok:       true,
			expected: "Module",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			version := version.New(newBuildInfoFunc(tc.module, tc.buildInfoSettings, tc.ok), "Dev")
			assert.Equal(t, version.Version(), tc.expected)
		})
	}
}

func newBuildInfoFunc(module debug.Module, buildSettings []debug.BuildSetting, ok bool) func() (*debug.BuildInfo, bool) {
	return func() (*debug.BuildInfo, bool) {
		buildInfo := debug.BuildInfo{
			Main:     module,
			Settings: buildSettings,
		}
		return &buildInfo, ok
	}
}
