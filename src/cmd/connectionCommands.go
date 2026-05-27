// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/connectionCommands.go
// Original timestamp: 2026/05/20 08:05:46

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"vmman4/connection"
)

var connCmd = &cobra.Command{
	Use:     "connection",
	Aliases: []string{"conn"},
	Short:   "Connection subcommands",
	Long:    `You need to provide one of the subcommands: ls, create, rm, info.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, create, or rm")
	},
}

var connLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List connection files",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection.ConnList(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var connRmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"del", "remove"},
	Short:   "Delete connection file(s)",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection.ConnRemove(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var connAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Short:   "Create an connection file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection.ConnCreate(); err != nil {
			fmt.Println(err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(connCmd)
	connCmd.AddCommand(connLsCmd, connRmCmd, connAddCmd)

	//connAddCmd.Flags().StringVar(&connection.ConnectionName, "connection", "", "Target hypervisor")
	connAddCmd.Flags().StringVar(&connection.ConnectionHost, "host", "", "Environment file.")
	connAddCmd.Flags().StringVar(&connection.ConnectionUser, "username", "", "Make vmman multi hypervisor-aware")
	connAddCmd.Flags().StringVar(&connection.ConnectionComment, "comments", "", "Make vmman multi hypervisor-aware")
	//// if one the above flag is set, the 3 others have to be set as well
	//connAddCmd.MarkFlagsRequiredTogether("host", "username", "comments")
}
