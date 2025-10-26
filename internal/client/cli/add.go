package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rookgm/gophkeeper/internal/models"
	"github.com/spf13/cobra"
	"os"
)

// TODO remove
func (c *secretCmd) Test(cmd *cobra.Command, args []string) error {
	fmt.Println("This is test, remove it")
	//secReq := models.SecretRequest{
	//	Name: "test_secret",
	//	Type: models.Text,
	//	Note: "test",
	//	Data: []byte("this is test data"),
	//}
	//_, err := c.secretSvc.CreateSecret(cmd.Context(), secReq, "123")

	id, _ := uuid.Parse("db08060a-d45a-4146-898d-be24bb4dddad")
	sec, err := c.secretSvc.GetSecret(cmd.Context(), id, "123")
	if err != nil {
		return fmt.Errorf("Error getting secret: %v\n", err)
	}

	fmt.Println(string(sec.Data))

	return err
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

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("Error marshalling binary data request: %v\n", err)
	}

	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Credential,
		Note: req.Note,
		Data: reqJSON,
	}

	resp, err := c.secretSvc.CreateSecret(cmd.Context(), secReq, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error creating secret: %v\n", err)
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

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("Error marshalling binary data request: %v\n", err)
	}

	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Text,
		Note: req.Note,
		Data: reqJSON,
	}

	resp, err := c.secretSvc.CreateSecret(cmd.Context(), secReq, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error creating secret: %v\n", err)
	}

	fmt.Printf("Successfully added text, ID: %s\n", resp.ID)

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

	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %v\n", fileName, err)
	}

	req := models.BinaryData{
		Name:     name,
		FileName: fileName,
		Note:     note,
		Data:     data,
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("Error marshalling binary data request: %v\n", err)
	}

	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Binary,
		Note: req.Note,
		Data: reqJSON,
	}

	resp, err := c.secretSvc.CreateSecret(cmd.Context(), secReq, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error creating secret: %v\n", err)
	}

	fmt.Printf("Successfully added binary data, ID: %s\n", resp.ID)

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

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("Error marshalling binary data request: %v\n", err)
	}

	secReq := models.SecretRequest{
		Name: req.Name,
		Type: models.Card,
		Note: req.Note,
		Data: reqJSON,
	}

	resp, err := c.secretSvc.CreateSecret(cmd.Context(), secReq, c.masterPassword)
	if err != nil {
		return fmt.Errorf("Error creating secret: %v\n", err)
	}

	fmt.Printf("Successfully added bank card, ID: %s\n", resp.ID)

	return nil
}
