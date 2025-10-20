package cli

import (
	"github.com/spf13/cobra"
)

func newVersionCmd(info BuildInfoPrinter) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			info.Print()
		},
	}
}
