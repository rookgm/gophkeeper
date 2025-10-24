package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd(userSvc UserService, secretSvc SecretService, buildInfo BuildInfoPrinter) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gophkeeper-client",
		Short: "A gophkeeper small password manager CLI",
		Long: `Manager passwords is a secure command-line password manager CLI that stores
	your secrets encrypted locally and remotely`,
	}

	rootCmd.AddCommand(
		newRegisterCmd(userSvc),
		newLoginCmd(userSvc),
		newSecretCmd(secretSvc),
		syncCmd,
		newVersionCmd(buildInfo),
	)
	return rootCmd
}
