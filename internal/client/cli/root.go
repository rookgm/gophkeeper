package cli

import (
	"fmt"
	"github.com/rookgm/gophkeeper/internal/build"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "gophkeeper-client",
	Short: "A gophkeeper small password manager CLI",
	Long: `Manager passwords is a secure command-line password manager CLI that stores
	your secrets encrypted locally and remotely`,
}

func Execute(info *build.AppBuildInfo) {
	rootCmd.AddCommand(
		registerCmd,
		loginCmd,
		secretCmd,
		syncCmd,
		versionCmd(info),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
