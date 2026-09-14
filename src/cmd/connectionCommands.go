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
	Use:     "conn",
	Aliases: []string{"connection"},
	Example: "vmman conn ls",
	Short:   "Connection subcommands",
	Long:    `You need to provide one of the subcommands: ls, create, rm, info.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, create, rm or info")
	},
}

var connLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Example: "vmman conn ls",
	Short:   "List saved connection files",
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection_mgt.ListConnections(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var connRmCmd = &cobra.Command{
	Use:     "rm <NAME...>",
	Aliases: []string{"del", "remove"},
	Example: "vmman conn rm myhypervisor",
	Short:   "Delete one or more connection files",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection_mgt.RemoveConnection(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var connAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create"},
	Example: "vmman conn add\nvmman conn add --host 192.168.1.10 --username root --comments \"lab hypervisor\"",
	Short:   "Interactively create a connection file",
	Long:    `Prompts for a connection name, host, username, and an optional comment, then saves the result to ~/.config/JFG/vmman4/<name>.json. Pass --host, --username, and/or --comments to pre-fill any of those and skip their prompt.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := connection_mgt.CreateConnection(); err != nil {
			fmt.Println(err.Error())
		}

	},
}

var connInfoCmd = &cobra.Command{
	Use:     "info <NAME...>",
	Aliases: []string{"explain"},
	Example: "vmman conn info myhypervisor\nvmman conn info myhypervisor1 myhypervisor2",
	Short:   "Print the details of one or more connection files",
	Long:    `You can list as many connection files as you wish, here`,
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
	connAddCmd.Flags().StringVar(&connection_mgt.ConnectionHost, "host", "", "Hostname or IP of the remote hypervisor (skips the host prompt)")
	connAddCmd.Flags().StringVar(&connection_mgt.ConnectionUser, "username", "", "SSH username for the remote hypervisor (skips the username prompt)")
	connAddCmd.Flags().StringVar(&connection_mgt.ConnectionComment, "comments", "", "Optional free-text comment for the connection (skips the comment prompt)")
	//// if one the above flag is set, the 3 others have to be set as well
	//connAddCmd.MarkFlagsRequiredTogether("host", "username", "comments")
}
