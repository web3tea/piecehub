package version

import (
	"fmt"
	"runtime/debug"
)

func GetVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		version := info.Main.Version
		var revision string
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				revision = setting.Value[:8]
				break
			}
		}

		if version == "(devel)" && revision != "" {
			return fmt.Sprintf("dev-%s", revision)
		}
		if version == "" {
			version = "unknown"
		}
		return version
	}
	return "unknown"
}
