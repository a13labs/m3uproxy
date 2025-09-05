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
	usersCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a user",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.PrintErrln("Usage: m3uproxy users remove <username>")
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

		err = auth.RemoveUser(args[0])
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		fmt.Println("User removed")
		os.Exit(0)
	},
}
