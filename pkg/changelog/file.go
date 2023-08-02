package changelog

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/flume/release-version/pkg/git"
)

var changelogFile = "CHANGELOG.md"

// Prepend prepend content to file
func Prepend(dir string, content string) (string, error) {
	filePath := path.Join(dir, changelogFile)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return changelogFile, fmt.Errorf("[Prepend] open file: %v", err)
	}
	defer f.Close()

	current, err := ioutil.ReadFile(filePath)
	if err != nil {
		return changelogFile, fmt.Errorf("[Prepend] read file: %v", err)
	}

	writer := bufio.NewWriter(f)
	_, _ = writer.WriteString(content)
	_, _ = writer.Write(current)

	if err := writer.Flush(); err != nil {
		return changelogFile, fmt.Errorf("[Prepend] flush file: %v", err)
	}

	return changelogFile, nil
}

// GitCommit adds CHANGELOG.md to Git
func GitCommit(dir string, version string, user *git.User) error {
	message := fmt.Sprintf("chore(release): %s", version)

	err := git.Commit(dir, changelogFile, message, user)
	if err != nil {
		return fmt.Errorf("[GitAdd] open repo: %v", err)
	}

	return nil
}
