package cli

import (
	"context"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/spf13/cobra"
)

type UserService interface {
	RegisterUser(ctx context.Context, user, password string) error
	LoginUser(ctx context.Context, user, password string) error
}

type SecretService interface {
	AddCredentials(ctx context.Context, req models.Credentials, masterPassword string) (*models.Credentials, error)
	AddText(ctx context.Context, req models.TextData, masterPassword string) (*models.TextData, error)
	AddBinary(ctx context.Context, req models.BinaryData, masterPassword string) (*models.BinaryData, error)
	AddBankCard(ctx context.Context, req models.BankCard, masterPassword string) (*models.BankCard, error)
}

type BuildInfoPrinter interface {
	Print()
}

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
