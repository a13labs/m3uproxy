package users

import (
	"os"

	"github.com/a13labs/m3uproxy/pkg/auth"
	"github.com/a13labs/m3uproxy/pkg/streamserver"
	rootCmd "github.com/a13labs/m3uproxy/server/cmd"

	"github.com/spf13/cobra"
)

func init() {
	usersCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new user",
	Run: func(cmd *cobra.Command, args []string) {
		// Add your code here to handle the "add" command
		if len(args) != 2 {
			cmd.PrintErrln("Usage: m3uproxy users add <username> <password>")
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

		err = auth.AddUser(args[0], args[1])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		cmd.Println("User added")
		os.Exit(0)
	},
}
