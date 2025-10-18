package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	RunE:  runRegister,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a user",
	RunE:  runLogin,
}

func init() {
	registerCmd.Flags().StringP("username", "u", "", "username to register with")
	loginCmd.Flags().StringP("username", "u", "", "username to login")
}

// readPassword reads password from terminal
func readPassword(msg string) (string, error) {
	fmt.Print(msg)
	pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("Error reading password: %v\n", err)
	}
	fmt.Println()

	return string(pwd), nil
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// runRegister registers a new user
//
// command: register
//
// flags: -u user_name
func runRegister(cmd *cobra.Command, args []string) error {
	// get username
	user, _ := cmd.Flags().GetString("username")
	if user == "" {
		return fmt.Errorf("username is required")
	}
	// read user password
	pwd, err := readPassword("password for " + user + ":")
	if err != nil {
		return fmt.Errorf("Error reading password: %v\n", err)
	}
	// TODO
	fmt.Println(pwd)

	return nil
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// runLogin login to a user
//
// commands: login
//
// flags: -u user_name
func runLogin(cmd *cobra.Command, args []string) error {
	// get username
	user, _ := cmd.Flags().GetString("username")
	if user == "" {
		return fmt.Errorf("username is required")
	}
	// read user password
	pwd, err := readPassword("password for " + user + ":")
	if err != nil {
		return fmt.Errorf("Error reading password: %v\n", err)
	}
	// TODO
	fmt.Println("password:", pwd)

	return nil
}
