package tnyc

import "github.com/maloquacious/semver"

var version = semver.Version{
	Major:      0,
	Minor:      1,
	Patch:      4,
	PreRelease: "alpha",
}

// Version returns the current version of tnyc.
func Version() semver.Version {
	return version
}
