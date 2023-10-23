package parser

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"regexp"
	"strings"
)

// SemVerChange describes the semver change type
type SemVerChange string

const (
	// Patch semver change
	Patch SemVerChange = "patch"
	// Minor semver change
	Minor SemVerChange = "minor"
	// Major semver change
	Major SemVerChange = "major"
)

// ConventionalCommit parsed commit
type ConventionalCommit struct {
	Hash         string
	Type         string
	Component    string
	Description  string
	Body         string
	Footer       string
	Breaking     string
	SemVerChange SemVerChange
	SemVer       string
}

// var pattern = regexp.MustCompile(`^(?:(\w+)\(?(\w+|\*)?\)?: (.+))(?:(?:\r?\n|$){0,2}(.+))?(?:(?:\r?\n|$){0,2}(.+))?(?:\r?\n|$){0,2}`)
var pattern = regexp.MustCompile(`^(?:(\w+)\(?(\w+|\*)?\)?: (.+))(?:(?:\r?\n|$){0,2}(.+\n)+)?(?:(?:\r?\n|$){0,2}(.+\n)+)?(?:\r?\n|$){0,2}`)
var versionPattern = regexp.MustCompile(`^((([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?)$`)
var breakingChange = regexp.MustCompile(`BREAKING\s?CHANGE:\s?([^\n]+)`)

// ParseCommits parses commits
func ParseCommits(dir string) ([]ConventionalCommit, error) {
	r, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("[ParseCommits] open repo: %v", err)
	}

	ref, err := r.Head()
	if err != nil {
		return nil, fmt.Errorf("[ParseCommits] head: %v", err)
	}

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("[ParseCommits] git log: %v", err)
	}

	var found = false
	var commits []ConventionalCommit

	tags := map[plumbing.Hash][]string{}
	iRefs, _ := r.Tags()
	_ = iRefs.ForEach(func(ref *plumbing.Reference) error {
		if strings.HasPrefix(ref.Strings()[0], "refs/tags/") {
			tags[ref.Hash()] = append(tags[ref.Hash()], strings.TrimPrefix(ref.Strings()[0], "refs/tags/"))
		}

		return nil
	})

	err = cIter.ForEach(func(c *object.Commit) error {
		if found {
			return nil
		}
		tmp := pattern.FindStringSubmatch(c.Message)

		// Skip commit that doesn't follow the conventional format
		if len(tmp) < 6 {
			return nil
		}

		commit := ConventionalCommit{
			Hash:         c.Hash.String(),
			Type:         tmp[1],
			Component:    tmp[2],
			Description:  tmp[3],
			Body:         tmp[4],
			Footer:       tmp[5],
			SemVerChange: Patch,
		}

		if commit.Component == "*" {
			commit.Component = ""
		}

		// Detect last semver bump from tag
		for _, tag := range tags[c.Hash] {
			matches := versionPattern.FindStringSubmatch(tag)
			if len(matches) > 0 {
				found = true
				commit.SemVer = matches[1]
				break
			}
		}

		// Detect last semver bump from commit
		tmp = versionPattern.FindStringSubmatch(commit.Description)

		if !found && commit.Type == "chore" && commit.Component == "release" && len(tmp) > 0 {
			found = true
			commit.SemVer = tmp[1]
		}

		if commit.Type == "feat" {
			commit.SemVerChange = Minor
		}

		if breakingChange.MatchString(commit.Body) {
			commit.SemVerChange = Major
			matches := breakingChange.FindAllStringSubmatch(commit.Body, -1)
			for _, m := range matches {
				commit.Breaking = commit.Breaking + m[len(m)-1] + "\n"
			}
		}

		if breakingChange.MatchString(commit.Footer) {
			commit.SemVerChange = Major
			matches := breakingChange.FindAllStringSubmatch(commit.Footer, -1)
			for _, m := range matches {
				commit.Breaking = commit.Breaking + m[len(m)-1] + "\n"
			}
		}

		commits = append(commits, commit)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("[ParseCommits] parse: %v", err)
	}

	return commits, nil
}

// ToSemVerChange converts string to SemVerChange
func ToSemVerChange(changeFlag string) (change SemVerChange) {
	switch changeFlag {
	case "patch":
		change = Patch
	case "minor":
		change = Minor
	case "major":
		change = Major
	}
	return change
}
