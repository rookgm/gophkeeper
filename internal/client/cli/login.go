package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

type loginCmd struct {
	userSvc UserService
}

func newLoginCmd(userSvc UserService) *cobra.Command {
	login := &loginCmd{userSvc: userSvc}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to a user",
		RunE:  login.login,
	}
	cmd.Flags().StringP("username", "u", "", "username to login")

	return cmd
}

// login performs login to a user
//
// commands: login
//
// flags: -u user_name
func (c *loginCmd) login(cmd *cobra.Command, args []string) error {
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
	// login user
	if err := c.userSvc.LoginUser(cmd.Context(), user, pwd); err != nil {
		return fmt.Errorf("error user login: %v\n", err)
	}

	fmt.Println("login successfully")

	return nil
}
