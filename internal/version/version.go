package version

import (
	"regexp"
	"runtime/debug"
)

type BuildInfoFunc func() (*debug.BuildInfo, bool)

type Version struct {
	readBuildInfo BuildInfoFunc
	version       string
}

func (v Version) Version() string {
	buildInfo, ok := v.readBuildInfo()
	if !ok {
		return v.version
	}
	ldFlagVersion, ok := v.ldflagVersion(buildInfo.Settings)
	if !ok {
		return buildInfo.Main.Version
	}
	return ldFlagVersion
}

func (v Version) ldflagVersion(settings []debug.BuildSetting) (string, bool) {
	regex := regexp.MustCompile(`.*main\.Version=(\S*)`)

	for _, s := range settings {
		if s.Key == "-ldflags" {
			matches := regex.FindStringSubmatch(s.Value)

			if len(matches) == 0 {
				return "", false
			}
			return matches[len(matches)-1], true
		}
	}
	return "", false
}

func New(readBuildInfo BuildInfoFunc, defaultVersion string) Version {
	return Version{
		readBuildInfo: readBuildInfo,
		version:       defaultVersion,
	}
}
