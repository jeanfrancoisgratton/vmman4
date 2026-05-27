// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/vmCommands.go
// Original timestamp: 2026/05/26 21:21:12

package cmd

import (
	"fmt"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
	"vmman4/inventory"
)

// vmCmd covers all the vm-related subcommands
var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Virtual machines management",
	Long:  `Commands to manage the VMs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.ScrollSign("Available commands: vm [ list | stop | start | restart ]"))
	},
}

var vmLsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all VMs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := inventory.VmInventory(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(vmCmd, vmLsCmd)
	vmCmd.AddCommand(vmLsCmd)
}
