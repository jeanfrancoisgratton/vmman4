// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/connectionCommands.go
// Original timestamp: 2026/05/20 08:05:46

package cmd

import (
	"fmt"

	"vmman4/connection_mgt"

	"github.com/spf13/cobra"
)

var connCmd = &cobra.Command{
	Use:     "connection_mgt",
	Aliases: []string{"conn"},
	Short:   "Connection subcommands",
	Long:    `You need to provide one of the subcommands: ls, create, rm, info.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, create, rm or info")
	},
}

var connLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List connection_mgt files",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection_mgt.ConnList(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var connRmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"del", "remove"},
	Short:   "Delete connection_mgt file(s)",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection_mgt.ConnRemove(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var connAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Short:   "Create an connection_mgt file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection_mgt.ConnCreate(); err != nil {
			fmt.Println(err.Error())
		}

	},
}

var connInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"explain"},
	Example: "vmman conn info FILE1[.json] FILE2[.json]... FILEn[.json]",
	Short:   "Prints the connection_mgt FILE[12n] information",
	Long:    `You can list as many connection_mgt files as you wish, here`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection_mgt.ExplainConnFile(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(connCmd)
	connCmd.AddCommand(connLsCmd, connRmCmd, connAddCmd, connInfoCmd)

	//connAddCmd.Flags().StringVar(&connection_mgt.ConnectionName, "connection_mgt", "", "Target hypervisor")
	connAddCmd.Flags().StringVar(&connection_mgt.ConnectionHost, "host", "", "Environment file.")
	connAddCmd.Flags().StringVar(&connection_mgt.ConnectionUser, "username", "", "Make vmman multi hypervisor-aware")
	connAddCmd.Flags().StringVar(&connection_mgt.ConnectionComment, "comments", "", "Make vmman multi hypervisor-aware")
	//// if one the above flag is set, the 3 others have to be set as well
	//connAddCmd.MarkFlagsRequiredTogether("host", "username", "comments")
}
