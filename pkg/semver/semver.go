package semver

import (
	"fmt"

	sv "github.com/coreos/go-semver/semver"
	"github.com/flume/release-version/pkg/parser"
)

// GetChange determinate semver changes (patch, minor, major)
func GetChange(commits []parser.ConventionalCommit) parser.SemVerChange {
	var change parser.SemVerChange = parser.Patch
	for _, commit := range commits {
		if change != parser.Major && commit.SemVerChange == parser.Minor {
			change = parser.Minor
		}
		if commit.SemVerChange == parser.Major {
			change = parser.Major
		}
	}
	return change
}

func GetLastVersion(commits []parser.ConventionalCommit) string {
	i := len(commits) - 1
	for i > -1 {
		commit := commits[i]
		if commit.SemVer != "" {
			return commit.SemVer
		}
		i--
	}

	return "0.0.0"
}

// GetVersion calculate version
func GetVersion(version string, change parser.SemVerChange) (string, error) {
	if version == "" {
		return "1.0.0", nil
	}

	v, err := sv.NewVersion(version)
	if err != nil {
		return "", fmt.Errorf(
			"[semver.GetVersion] parse version (%s): %v",
			version,
			err,
		)
	}

	switch change {
	case parser.Patch:
		v.BumpPatch()
	case parser.Minor:
		v.BumpMinor()
	case parser.Major:
		v.BumpMajor()
	default:
		return "", fmt.Errorf("Invalid change type %s", change)
	}

	return v.String(), nil
}
