package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/flume/release-version/pkg/parser"
	"github.com/spf13/cobra"
)

// GetParseCmd returns the parse cmd
func GetParseCmd() *cobra.Command {
	var branch string
	// Default dir is the working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cmdParse = &cobra.Command{
		Use:   "parse",
		Short: "Parses commits",
		Long:  `parsing conventional commits.`,
		Run: func(cmd *cobra.Command, args []string) {
			commits, err := parser.ParseCommits(dir)

			if err != nil {
				panic(err)
			}
			b, _ := json.MarshalIndent(commits, "", "  ")
			fmt.Println(string(b))
		},
	}

	cmdParse.Flags().StringVarP(
		&dir,
		"repository",
		"r",
		"",
		"Repository directory",
	)

	cmdParse.Flags().StringVarP(
		&branch,
		"branch",
		"b",
		"",
		"Which branch to write run against",
	)
	return cmdParse
}
