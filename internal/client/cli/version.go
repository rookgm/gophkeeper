package cli

import (
	"github.com/rookgm/gophkeeper/internal/build"
	"github.com/spf13/cobra"
)

func versionCmd(info *build.AppBuildInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			info.Print()
		},
	}
}
