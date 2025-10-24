package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/spf13/cobra"
)

type SecretService interface {
	CreateSecret(ctx context.Context, req models.SecretRequest, masterPassword string) (*models.SecretResponse, error)
	GetSecret(ctx context.Context, id uuid.UUID, masterPassword string) (*models.SecretResponse, error)
	DeleteSecret(ctx context.Context, id uuid.UUID, masterPassword string) (*models.SecretResponse, error)
	UpdateSecret(ctx context.Context, id uuid.UUID, req models.SecretRequest, masterPassword string) (*models.SecretResponse, error)
}

type secretCmd struct {
	secretSvc      SecretService
	masterPassword string
}

func newSecretCmd(secretSvc SecretService) *cobra.Command {
	sec := &secretCmd{secretSvc: secretSvc}
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "Manage secrets",
	}
	// create secret adder command
	secretAddCmd := &cobra.Command{Use: "add", Short: "Add a secret"}
	// add commands: add, delete, get to main secret cmd
	cmd.AddCommand(secretAddCmd)
	cmd.AddCommand(&cobra.Command{Use: "get", Short: "Get a secret", Args: cobra.ExactArgs(1), PreRunE: sec.auth, RunE: sec.runSecretGetCmd})
	cmd.AddCommand(&cobra.Command{Use: "delete", Short: "Delete a secret", Args: cobra.ExactArgs(1), PreRunE: sec.auth, RunE: sec.runSecretDeleteCmd})

	addCredentialsCmd := &cobra.Command{Use: "credentials", Short: "Add credentials", PreRunE: sec.auth, RunE: sec.runAddCredentials}
	addTextCmd := &cobra.Command{Use: "text", Short: "Add text", PreRunE: sec.auth, RunE: sec.runAddTextCmd}
	addBinaryCmd := &cobra.Command{Use: "binary", Short: "Add binary data", PreRunE: sec.auth, RunE: sec.runAddBinaryCmd}
	addBankCardCmd := &cobra.Command{Use: "card", Short: "Add bank card", PreRunE: sec.auth, RunE: sec.runAddBankCardCmd}

	// add command to secret adder command
	secretAddCmd.AddCommand(addCredentialsCmd, addTextCmd, addBinaryCmd, addBankCardCmd)

	// common flags
	secretAddCmd.PersistentFlags().StringP("name", "n", "", "secret name")
	secretAddCmd.PersistentFlags().StringP("note", "e", "", "secret note")
	// credentials
	addCredentialsCmd.Flags().StringP("login", "l", "", "login")
	// text
	addTextCmd.Flags().StringP("content", "c", "", "text content")
	// binary data from file
	addBinaryCmd.Flags().StringP("filename", "p", "", "file name")
	// bank card
	addBankCardCmd.Flags().String("number", "", "full card number")
	addBankCardCmd.Flags().String("expmonth", "", "two-digit expiration month of the card")
	addBankCardCmd.Flags().String("expyear", "", "four-digit expiration year of the card")
	addBankCardCmd.Flags().String("holdername", "", "card holder name")
	addBankCardCmd.Flags().String("address", "", "an object containing details of the cardholder's billing address")
	addBankCardCmd.Flags().String("type", "", "type of credit card (e.g., Visa, Mastercard, American Express)")
	addBankCardCmd.Flags().String("issue", "", "name of the bank that issued the bank card")

	return cmd
}

func (c *secretCmd) auth(cmd *cobra.Command, args []string) error {
	master, err := readPassword("Enter master password: ")
	if err != nil {
		return fmt.Errorf("Error reading master password: %v\n", err)
	}
	c.masterPassword = master

	return nil
}

// runSecretGetCmd gets secret info
//
// command: secret get <secret_id>
func (c *secretCmd) runSecretGetCmd(cmd *cobra.Command, args []string) error {
	id, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("Error parsing secret id: %v\n", err)
	}
	sec, err := c.secretSvc.GetSecret(cmd.Context(), id, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error getting secret: %v\n", err)
	}

	fmt.Printf("ID: %s\n", sec.ID)
	fmt.Printf("Type: %s\n", sec.Type.String())
	fmt.Printf("Name: %s\n", sec.Name)
	if sec.Note != "" {
		fmt.Printf("Note: %s\n", sec.Note)
	}
	fmt.Printf("Created: %s\n", sec.CreatedAt)
	fmt.Printf("Updated: %s\n", sec.UpdatedAt)

	data, err := json.MarshalIndent(sec.Data, "", " ")
	if err != nil {
		return fmt.Errorf("Error marshalling secret data: %v\n", err)
	}
	fmt.Printf("Data: %s\n", string(data))

	return nil
}

// runSecretDeleteCmd removes secret
//
// command: secret delete <secret_id>
func (c *secretCmd) runSecretDeleteCmd(cmd *cobra.Command, args []string) error {
	id, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("Error parsing secret id: %v\n", err)
	}

	var resp string
	fmt.Printf("Are you sure you want to delete the secret %s? (y/N)", id)
	fmt.Scanln(&resp)

	if resp != "y" && resp != "Y" {
		return nil
	}
	if _, err := c.secretSvc.DeleteSecret(cmd.Context(), id, c.masterPassword); err != nil {
		return fmt.Errorf("Error deleting secret: %v\n", err)
	}

	return nil
}
