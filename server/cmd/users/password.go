/*
Copyright © 2024 Alexandre Pires

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
		err := streamserver.LoadServerConfig(rootCmd.ConfigFile)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		err = auth.InitializeAuth(streamserver.Config.Auth)
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
