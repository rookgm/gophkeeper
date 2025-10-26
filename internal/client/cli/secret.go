package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
)

type SecretService interface {
	CreateSecret(ctx context.Context, req models.SecretRequest, masterPassword string) (*models.SecretResponse, error)
	GetSecret(ctx context.Context, id uuid.UUID, masterPassword string) (*models.SecretResponse, error)
	DeleteSecret(ctx context.Context, id uuid.UUID) error
	UpdateSecret(ctx context.Context, id uuid.UUID, req models.SecretRequest, masterPassword string) error
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
	// add commands: add, delete, get, update to main secret cmd
	cmd.AddCommand(secretAddCmd)
	cmd.AddCommand(&cobra.Command{Use: "get", Short: "Get a secret", Args: cobra.ExactArgs(1), PreRunE: sec.auth, RunE: sec.runSecretGetCmd})
	cmd.AddCommand(&cobra.Command{Use: "delete", Short: "Delete a secret", Args: cobra.ExactArgs(1), RunE: sec.runSecretDeleteCmd})
	cmd.AddCommand(&cobra.Command{Use: "update", Short: "Update a secret", Args: cobra.ExactArgs(1), PreRunE: sec.auth, RunE: sec.runSecretUpdateCmd})

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
	// print common secret information
	header := fmt.Sprintf("===%s details===", sec.Type.String())
	fmt.Println(strings.ToUpper(header))
	fmt.Printf("ID: %s\n", sec.ID)
	fmt.Printf("Name: %s\n", sec.Name)
	if sec.Note != "" {
		fmt.Printf("Note: %s\n", sec.Note)
	}
	// print specific secret information
	switch sec.Type {
	case models.Credential:
		var cred models.Credentials
		err = json.Unmarshal(sec.Data, &cred)
		if err != nil {
			return fmt.Errorf("Error unmarshalling credentials: %v\n", err)
		}
		fmt.Printf("Login: %v\n", cred.Login)
		fmt.Printf("Password: %v\n", cred.Password)
	case models.Text:
		var textData models.TextData
		err = json.Unmarshal(sec.Data, &textData)
		if err != nil {
			return fmt.Errorf("Error unmarshalling text data: %v\n", err)
		}
		fmt.Printf("Content: %v\n", textData)
	case models.Binary:
		var binaryData models.BinaryData
		err = json.Unmarshal(sec.Data, &binaryData)
		if err != nil {
			return fmt.Errorf("Error unmarshalling binary data: %v\n", err)
		}
		// write file to current directory
		err := os.WriteFile(binaryData.FileName, binaryData.Data, 0600)
		if err != nil {
			return fmt.Errorf("Error writing %s: %v\n", binaryData.FileName, err)
		}
		fmt.Printf("File %s has been saved\n", binaryData.FileName)
	case models.Card:
		var card models.BankCard
		err = json.Unmarshal(sec.Data, &card)
		if err != nil {
			return fmt.Errorf("Error unmarshalling card data: %v\n", err)
		}
		fmt.Printf("Number: %v\n", card.CardNumber)
		fmt.Printf("Expiration Month: %v\n", card.ExpirationMonth)
		fmt.Printf("Expiration Year: %v\n", card.ExpirationYear)
		fmt.Printf("Holder name: %v\n", card.CardHolderName)
		fmt.Printf("CVV: %v\n", card.Cvv)
		fmt.Printf("Billing address: %v\n", card.BillingAddress)
		fmt.Printf("Type: %v\n", card.CardType)
		fmt.Printf("Issuing Bank: %v\n", card.IssuingBank)

	default:
		return fmt.Errorf("Invalid secret type: %s\n", sec.Type)
	}

	fmt.Printf("Created: %s\n", sec.CreatedAt.Format(time.DateTime))
	fmt.Printf("Updated: %s\n", sec.UpdatedAt.Format(time.DateTime))

	return nil
}

// runSecretDeleteCmd removes secret
//
// command: secret update <secret_id>
func (c *secretCmd) runSecretDeleteCmd(cmd *cobra.Command, args []string) error {
	id, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("Error parsing secret id: %v\n", err)
	}

	if err := c.secretSvc.DeleteSecret(cmd.Context(), id); err != nil {
		return fmt.Errorf("Error deleting secret: %v\n", err)
	}

	fmt.Println("Secret has been successfully deleted")

	return nil
}

// runSecretUpdateCmd update secret
//
// command: secret update <secret_id>
func (c *secretCmd) runSecretUpdateCmd(cmd *cobra.Command, args []string) error {
	id, err := uuid.Parse(args[0])
	if err != nil {
		return fmt.Errorf("Error parsing secret id: %v\n", err)
	}

	req := models.SecretRequest{}

	err = c.secretSvc.UpdateSecret(cmd.Context(), id, req, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error updating secret: %v\n", err)
	}

	fmt.Println("Secret has been successfully updated")

	return nil
}
