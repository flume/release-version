package changelog

import (
	"fmt"

	"github.com/flume/release-version/pkg/git"
	"github.com/flume/release-version/pkg/parser"
)

// Save generates and adds changelog.md to Git
func Save(
	dir string,
	newVersion string,
	lastVersion string,
	change parser.SemVerChange,
	commits []parser.ConventionalCommit,
	user *git.User,
	branch string,
) (
	string,
	string,
	error,
) {
	// get a remote path to try to use in the markdown
	rPath, _ := git.GetRemotePath(dir)

	// Generate changelog
	markdown := Generate(newVersion, lastVersion, change, commits, rPath)

	// Write changelog
	file, err := Prepend(dir, markdown)
	if err != nil {
		return file, markdown, fmt.Errorf("[Save] prepend: %v", err)
	}

	// Add to Git
	err = GitCommit(dir, newVersion, user, branch)
	if err != nil {
		return file, markdown, fmt.Errorf("[Save] git commit: %v", err)
	}

	return file, markdown, nil
}
