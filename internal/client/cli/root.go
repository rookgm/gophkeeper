package cli

import (
	"context"
	"github.com/spf13/cobra"
)

type UserService interface {
	RegisterUser(ctx context.Context, user, password string) error
	LoginUser(ctx context.Context, user, password string) error
}

type BuildInfoPrinter interface {
	Print()
}

func NewRootCmd(userSvc UserService, buildInfo BuildInfoPrinter) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gophkeeper-client",
		Short: "A gophkeeper small password manager CLI",
		Long: `Manager passwords is a secure command-line password manager CLI that stores
	your secrets encrypted locally and remotely`,
	}

	rootCmd.AddCommand(
		newRegisterCmd(userSvc),
		newLoginCmd(userSvc),
		secretCmd,
		syncCmd,
		newVersionCmd(buildInfo),
	)
	return rootCmd
}
