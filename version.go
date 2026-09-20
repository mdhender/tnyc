package tnyc

import "github.com/maloquacious/semver"

var (
	version = semver.Version{
		Major:      0,
		Minor:      1,
		Patch:      6,
		PreRelease: "alpha",

		// Automatically populate build metadata with commit info
		Build: semver.Commit(), // Uses Git commit hash from build info
	}
)

func Version() semver.Version {
	return version
}
