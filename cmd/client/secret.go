package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var secretCmd = &cobra.Command{
	Use:    "secret",
	Short:  "Manage secrets",
	PreRun: auth,
}

var secretAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a secret",
}

var addCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "Add credentials",
	RunE:  runAddCredentials,
}

var addTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Add text",
	RunE:  runAddTextCmd,
}

var addBinaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Add binary data",
	RunE:  runAddBinaryCmd,
}

var addBankCardCmd = &cobra.Command{
	Use:   "card",
	Short: "Add bank card",
	RunE:  runAddBankCardCmd,
}

var secretGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runSecretGetCmd,
}

var secretDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a secret",
	Args:  cobra.ExactArgs(1),
	RunE:  runSecretDeleteCmd,
}

func init() {
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

	secretCmd.AddCommand(secretAddCmd, secretGetCmd, secretDeleteCmd)
	secretAddCmd.AddCommand(addCredentialsCmd, addTextCmd, addBinaryCmd, addBankCardCmd)

}

// TODO
func auth(cmd *cobra.Command, args []string) {

}

// runAddCredentials adds credentials data
//
// command: secret add credentials
//
// flags
// -n "name"
// -l "login"
// -e "note"
func runAddCredentials(cmd *cobra.Command, args []string) error {
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
	// TODO
	fmt.Println(pwd)

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
func runAddTextCmd(cmd *cobra.Command, args []string) error {
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
func runAddBinaryCmd(cmd *cobra.Command, args []string) error {
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

	fmt.Println(name, fileName, note)

	return nil
}

// runAddBankCardCmd adds bank card data
//
// command: secret add card
//
// flags:
// -n "name"
// -e "note"
func runAddBankCardCmd(cmd *cobra.Command, args []string) error {
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
	// TODO
	fmt.Println("cvv:", cvv)

	return nil
}

// runSecretGetCmd gets secret info
//
// command: secret get <secret_id>
func runSecretGetCmd(cmd *cobra.Command, args []string) error {
	id := args[0]
	fmt.Println(id)
	return nil
}

// runSecretDeleteCmd removes secret
//
// command: secret delete <secret_id>
func runSecretDeleteCmd(cmd *cobra.Command, args []string) error {
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
