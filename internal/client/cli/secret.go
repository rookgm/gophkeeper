package cli

import (
	"errors"
	"fmt"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/spf13/cobra"
	"os"
)

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

// runAddCredentials adds credentials data
//
// command: secret add credentials
//
// flags
// -n "name"
// -l "login"
// -e "note"
func (c *secretCmd) runAddCredentials(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		fmt.Print("Name: ")
		fmt.Scanln(&name)
	}

	login, _ := cmd.Flags().GetString("login")
	if login == "" {
		fmt.Print("Login: ")
		fmt.Scanln(&login)
	}

	note, _ := cmd.Flags().GetString("note")
	if note == "" {
		fmt.Print("Note: ")
		fmt.Scanln(&note)
	}

	// read password for login
	pwd, err := readPassword("password for " + login + ":")
	if err != nil {
		return fmt.Errorf("Error reading password: %v\n", err)
	}

	req := models.Credentials{
		Name:     name,
		Note:     note,
		Login:    login,
		Password: pwd,
	}

	resp, err := c.secretSvc.AddCredentials(cmd.Context(), req, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error adding credentials: %v\n", err)
	}

	fmt.Printf("Successfully added credentials, ID: %s\n", resp.ID)

	return nil
}

// runAddTextCmd adds text data
//
// command: secret add text
//
// flags
// -n "name"
// -c "content"
// -e "note"
func (c *secretCmd) runAddTextCmd(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		fmt.Print("Name: ")
		fmt.Scanln(&name)
	}

	content, _ := cmd.Flags().GetString("content")
	if content == "" {
		fmt.Print("Content: ")
		fmt.Scan(&content)
	}

	note, _ := cmd.Flags().GetString("note")
	if note == "" {
		fmt.Print("Note: ")
		fmt.Scanln(&note)
	}

	req := models.TextData{
		Name:    name,
		Note:    note,
		Content: content,
	}

	res, err := c.secretSvc.AddText(cmd.Context(), req, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error adding text: %v\n", err)
	}

	fmt.Printf("Successfully added text, ID: %s\n", res.ID)

	return nil
}

// runAddBinaryCmd adds binary data
//
// command: secret add binary
//
// flags
// -n "name"
// -p "file_name"
// -e "note"
func (c *secretCmd) runAddBinaryCmd(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		fmt.Print("Name: ")
		fmt.Scanln(&name)
	}

	fileName, _ := cmd.Flags().GetString("filename")
	if fileName == "" {
		fmt.Print("Filename: ")
		fmt.Scanln(&fileName)
	}

	note, _ := cmd.Flags().GetString("note")
	if note == "" {
		fmt.Print("Note: ")
		fmt.Scanln(&note)
	}

	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("File %s does not exist\n", fileName)
	}

	req := models.BinaryData{
		Name:     name,
		FileName: fileName,
		Note:     note,
		Data:     nil,
	}

	res, err := c.secretSvc.AddBinary(cmd.Context(), req, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error adding binary: %v\n", err)
	}

	fmt.Printf("Successfully added binary, ID: %s\n", res.ID)

	return nil
}

// runAddBankCardCmd adds bank card data
//
// command: secret add card
//
// flags:
// -n "name"
// -e "note"
func (c *secretCmd) runAddBankCardCmd(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	if name == "" {
		fmt.Print("Name: ")
		fmt.Scanln(&name)
	}

	note, _ := cmd.Flags().GetString("note")
	if note == "" {
		fmt.Print("Note: ")
		fmt.Scanln(&note)
	}

	number, _ := cmd.Flags().GetString("number")
	if number == "" {
		fmt.Print("Card number: ")
		fmt.Scanln(&number)
	}

	expmonth, _ := cmd.Flags().GetString("expmonth")
	if expmonth == "" {
		fmt.Print("Card expiration month: ")
		fmt.Scanln(&expmonth)
	}

	expyear, _ := cmd.Flags().GetString("expyear")
	if expyear == "" {
		fmt.Print("Card expiration year: ")
		fmt.Scanln(&expyear)
	}

	holdername, _ := cmd.Flags().GetString("holdername")
	if holdername == "" {
		fmt.Print("Card holder name: ")
		fmt.Scanln(&holdername)
	}

	address, _ := cmd.Flags().GetString("address")
	if address == "" {
		fmt.Print("Cardholder's billing address: ")
		fmt.Scanln(&address)
	}

	cardType, _ := cmd.Flags().GetString("type")
	if cardType == "" {
		fmt.Print("Card type: ")
		fmt.Scanln(&cardType)
	}

	issue, _ := cmd.Flags().GetString("issue")
	if issue == "" {
		fmt.Print("Issue name: ")
		fmt.Scanln(&issue)
	}

	// read ccv
	cvv, err := readPassword("CCV: ")
	if err != nil {
		return fmt.Errorf("Error reading CVV: %v\n", err)
	}

	req := models.BankCard{
		Name:            name,
		Note:            note,
		CardNumber:      number,
		ExpirationMonth: expmonth,
		ExpirationYear:  expyear,
		CardHolderName:  holdername,
		Cvv:             cvv,
		BillingAddress:  address,
		CardType:        cardType,
		IssuingBank:     issue,
	}

	resp, err := c.secretSvc.AddBankCard(cmd.Context(), req, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error adding bank card: %v\n", err)
	}

	fmt.Printf("Successfully added bank card, ID: %s\n", resp.ID)

	return nil
}

// runSecretGetCmd gets secret info
//
// command: secret get <secret_id>
func (c *secretCmd) runSecretGetCmd(cmd *cobra.Command, args []string) error {
	id := args[0]
	fmt.Println(id)
	return nil
}

// runSecretDeleteCmd removes secret
//
// command: secret delete <secret_id>
func (c *secretCmd) runSecretDeleteCmd(cmd *cobra.Command, args []string) error {
	id := args[0]
	fmt.Println(id)

	var resp string
	fmt.Printf("Are you sure you want to delete the secret %s? (y/N)", id)
	fmt.Scanln(&resp)

	if resp != "y" && resp != "Y" {
		return nil
	}

	return nil
}
