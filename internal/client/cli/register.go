package cli

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type UserService interface {
	RegisterUser(ctx context.Context, user, password string) error
	LoginUser(ctx context.Context, user, password string) error
}

type regCmd struct {
	userSvc UserService
}

func newRegisterCmd(userSvc UserService) *cobra.Command {
	reg := &regCmd{userSvc: userSvc}
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new user",
		RunE:  reg.register,
	}
	cmd.Flags().StringP("username", "u", "", "username to register with")

	return cmd
}

// register registers a new user
//
// command: register
//
// flags: -u user_name
func (c *regCmd) register(cmd *cobra.Command, args []string) error {
	// get username
	user, _ := cmd.Flags().GetString("username")
	if user == "" {
		return fmt.Errorf("username is required")
	}
	// read user password
	pwd, err := readPassword("password for " + user + ":")
	if err != nil {
		return fmt.Errorf("error reading password: %v\n", err)
	}

	if err := c.userSvc.RegisterUser(cmd.Context(), user, pwd); err != nil {
		return fmt.Errorf("error registering user: %v\n", err)
	}

	fmt.Printf("user: %s registered successfully\n", user)

	return nil
}

// readPassword reads password from terminal
func readPassword(msg string) (string, error) {
	fmt.Print(msg)
	pwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("error reading password: %v\n", err)
	}
	fmt.Println()

	return string(pwd), nil
}
