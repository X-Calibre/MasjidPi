package updates

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var installedVersionPattern = regexp.MustCompile(
	`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-rc\.([1-9][0-9]*))?$`,
)

type installedVersion struct {
	StableVersion
	ReleaseCandidate bool
	RCNumber         uint64
}

func parseInstalledVersion(
	value string,
) (installedVersion, error) {
	normalized := strings.TrimSuffix(value, "-image")
	match := installedVersionPattern.FindStringSubmatch(
		normalized,
	)
	if match == nil {
		return installedVersion{}, fmt.Errorf(
			"updates: %q is not a supported installed version",
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
			return installedVersion{}, fmt.Errorf(
				"updates: parse installed version %q: %w",
				value,
				err,
			)
		}
		parts[index] = parsed
	}

	installed := installedVersion{
		StableVersion: StableVersion{
			Major: parts[0],
			Minor: parts[1],
			Patch: parts[2],
		},
	}

	if match[4] != "" {
		rcNumber, err := strconv.ParseUint(
			match[4],
			10,
			64,
		)
		if err != nil {
			return installedVersion{}, fmt.Errorf(
				"updates: parse installed RC %q: %w",
				value,
				err,
			)
		}
		installed.ReleaseCandidate = true
		installed.RCNumber = rcNumber
	}

	return installed, nil
}

// IsNewerThanInstalled reports whether a stable release should be offered to
// the running appliance. A stable version supersedes an RC of the same core
// version.
func IsNewerThanInstalled(
	candidate StableVersion,
	current string,
) (bool, error) {
	installed, err := parseInstalledVersion(current)
	if err != nil {
		return false, err
	}

	comparison := candidate.Compare(
		installed.StableVersion,
	)
	if comparison != 0 {
		return comparison > 0, nil
	}

	return installed.ReleaseCandidate, nil
}
