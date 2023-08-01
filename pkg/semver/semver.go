package semver

import (
	"fmt"
	"github.com/flume/release-version/pkg/git"

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

func GetLastVersion(dir string) (string, error) {
	tag, err := git.GetLatestTag(dir)
	if err != nil {
		return "", fmt.Errorf("[semver.GetLastVersion] get latest tag: %v", err)
	}

	if tag == "" {
		return "0.0.0", nil
	}

	return tag, nil
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
