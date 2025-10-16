package cmd

import "github.com/spf13/cobra"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run:   runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {

}
