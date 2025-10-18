package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run:   runVersion,
}

// printBuildInfo prints application build info to stdout
func printBuildInfo() {
	fmt.Printf(
		"Build version: %s\n"+
			"Build date: %s\n"+
			"Build commit: %s\n",
		BuildVersion,
		BuildDate,
		BuildCommit)
}

// runVersion prints client build info to stdout
func runVersion(cmd *cobra.Command, args []string) {
	printBuildInfo()
}
