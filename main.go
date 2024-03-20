package main

import (
	"github.com/spf13/cobra"

	"github.com/flume/release-version/cmd"
)

func main() {
	var rootCmd = &cobra.Command{Use: "unchain"}
	rootCmd.AddCommand(cmd.GetParseCmd())
	rootCmd.AddCommand(cmd.GetReleaseCmd())
	rootCmd.AddCommand(cmd.GetSemverCmd())
	_ = rootCmd.Execute()
}
