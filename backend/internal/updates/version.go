package updates

import (
	"fmt"
	"regexp"
	"strconv"
)

var stableVersionPattern = regexp.MustCompile(
	`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`,
)

type StableVersion struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func ParseStableVersion(value string) (StableVersion, error) {
	match := stableVersionPattern.FindStringSubmatch(value)
	if match == nil {
		return StableVersion{}, fmt.Errorf(
			"updates: %q is not a stable semantic version",
			value,
		)
	}

	parts := [3]uint64{}
	for index := range parts {
		parsed, err := strconv.ParseUint(
			match[index+1],
			10,
			64,
		)
		if err != nil {
			return StableVersion{}, fmt.Errorf(
				"updates: parse version %q: %w",
				value,
				err,
			)
		}
		parts[index] = parsed
	}

	return StableVersion{
		Major: parts[0],
		Minor: parts[1],
		Patch: parts[2],
	}, nil
}

// Compare returns -1 when v is older, 0 when equal, and 1 when newer.
func (v StableVersion) Compare(other StableVersion) int {
	left := [...]uint64{v.Major, v.Minor, v.Patch}
	right := [...]uint64{
		other.Major,
		other.Minor,
		other.Patch,
	}

	for index := range left {
		switch {
		case left[index] < right[index]:
			return -1
		case left[index] > right[index]:
			return 1
		}
	}

	return 0
}
