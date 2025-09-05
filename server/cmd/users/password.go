package users

import (
	"fmt"
	"os"

	"github.com/a13labs/m3uproxy/pkg/auth"
	"github.com/a13labs/m3uproxy/pkg/streamserver"
	rootCmd "github.com/a13labs/m3uproxy/server/cmd"

	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(passwordCmd)
}

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Change a user password",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.PrintErrln("Usage: m3uproxy users password <username> <password>")
			os.Exit(1)
		}
		c := streamserver.NewServerConfig(rootCmd.ConfigFile)

		if c == nil {
			cmd.PrintErrln("Error loading config")
			os.Exit(1)
		}

		err := auth.InitializeAuth(c.Get().Auth)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		err = auth.ChangePassword(args[0], args[1])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		fmt.Println("Password changed")
		os.Exit(0)
	},
}
