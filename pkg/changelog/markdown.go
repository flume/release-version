package changelog

import (
	"bytes"
	"fmt"
	"time"

	"github.com/flume/release-version/pkg/parser"
)

// Generate generate markdown output
func Generate(version, lastVersion string, change parser.SemVerChange, commits []parser.ConventionalCommit, remotePath string) string {
	var out bytes.Buffer
	var patch = false
	var minor = false
	var major = false

	// Tag Header
	date := time.Now().Format("2006-01-02")

	headerLevel := "###"
	switch change {
	case parser.Major:
		headerLevel = "#"
	case parser.Minor:
		headerLevel = "##"
	}

	out.WriteString(
		fmt.Sprintf("%s [%s](%v) (%s)\n\n",
			headerLevel, version, getTagComparisonUrls(remotePath, lastVersion, version), date))

	// Patch
	for _, commit := range commits {
		if commit.SemVerChange == parser.Patch &&
			// Skip non user facing commits from changelog
			commit.Type != "test" && commit.Type != "chore" && commit.Type != "refactor" {

			if !patch {
				out.WriteString("### Bug Fixes\n")
			}
			out.WriteString(getCommitLine(&commit, remotePath))
			patch = true
		}
	}
	if patch {
		out.WriteString("\n")
	}

	// Minor
	for _, commit := range commits {
		if commit.SemVerChange == parser.Minor {
			if !minor {
				out.WriteString("\n### Features\n")
			}
			out.WriteString(getCommitLine(&commit, remotePath))
			minor = true
		}
	}
	if minor {
		out.WriteString("\n")
	}

	// Major
	for _, commit := range commits {
		if commit.SemVerChange == parser.Major {
			if !major {
				out.WriteString("\n### Breaking Changes\n")
			}
			out.WriteString(getBreakingLine(&commit))
			major = true
		}
	}
	if major {
		out.WriteString("\n")
	}

	// No user facing commit
	if !patch && !minor && !major {
		out.WriteString("* There is no user facing commit in this version\n")
	}

	out.WriteString("\n\n")

	return out.String()
}

func getCommitLine(commit *parser.ConventionalCommit, remotePath string) string {
	var out bytes.Buffer

	out.WriteString("\n* ")
	if len(commit.Component) > 0 {
		c := fmt.Sprintf("**%s:** ", commit.Component)
		out.WriteString(c)
	}
	out.WriteString(commit.Description)
	out.WriteString(" ")
	out.WriteString(fmt.Sprintf("([%v](%v))", commit.Hash[:7], getCommitUrl(remotePath, commit.Hash)))

	return out.String()
}

func getBreakingLine(commit *parser.ConventionalCommit) string {
	var out bytes.Buffer

	out.WriteString("\n* ")
	out.WriteString(commit.Breaking)
	out.WriteString(" ")
	out.WriteString(commit.Hash)

	return out.String()
}

// getCommitUrl is tailored to github
func getCommitUrl(remotePath, commitHash string) string {
	return fmt.Sprintf("https://%v/commit/%v", remotePath, commitHash)
}

// getTagComparisonUrls is tailored to github
func getTagComparisonUrls(remotePath, lastVersion, newVersion string) string {
	return fmt.Sprintf("https://%v/compare/%v...%v", remotePath, lastVersion, newVersion)
}
