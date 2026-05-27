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
	"vmman4/vmmanagement"
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

var vmStartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"up"},
	Short:   "Start one or many VMs",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vmmanagement.Start(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStartAllCmd = &cobra.Command{
	Use:   "startall",
	Short: "Start all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vmmanagement.StartAll(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(vmCmd, vmLsCmd, vmStartCmd, vmStartAllCmd)
	vmCmd.AddCommand(vmLsCmd, vmStartCmd, vmStartAllCmd)
}
