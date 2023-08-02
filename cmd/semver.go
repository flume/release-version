package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/flume/release-version/pkg/parser"
	"github.com/flume/release-version/pkg/semver"
	"github.com/spf13/cobra"
)

// GetSemverCmd returns the semver cmd
func GetSemverCmd() *cobra.Command {
	// Default dir is the working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cmdSemver = &cobra.Command{
		Use:   "semver",
		Short: "Next semver version",
		Long:  `semver is calculating the next SemVer version.`,
		Run: func(cmd *cobra.Command, args []string) {
			commits, err := parser.ParseCommits(dir)

			if err != nil {
				panic(err)
			}

			change := semver.GetChange(commits)
			fmt.Printf("Change Detected: %v\n", change)

			ver, _ := semver.GetVersion(semver.GetLastVersion(commits), change)
			fmt.Printf("Next Version: %v\n", ver)
		},
	}

	cmdSemver.Flags().StringVarP(
		&dir,
		"repository",
		"r",
		"",
		"Repository directory",
	)

	return cmdSemver
}
