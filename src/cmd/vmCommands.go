// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/vmCommands.go
// Original timestamp: 2026/05/26 21:21:12

package cmd

import (
	"fmt"

	"vmman4/vm_management"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
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
		if err := vm_management.VmInventory(); err != nil {
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
		if err := vm_management.Start(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStartAllCmd = &cobra.Command{
	Use:   "startall",
	Short: "Start all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_management.StartAll(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"down"},
	Short:   "Stop one or many VMs",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_management.Stop(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStopAllCmd = &cobra.Command{
	Use:   "stopall",
	Short: "Stop all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_management.StopAll(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmConsoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open a console session on the VM",
	Long:  "a -f flag will force the connection to the console, if that connection was already opened.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_management.Console(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmRenameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a VM",
	Long: `Be aware that that VM's underlying disk won't change name.
Also, if the VM holds any snapshot, they need to be removed before effecting the rename`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_management.Rename(args[0], args[1]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmRemoveCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "destroy", "delete"},
	Short:   "Remove one or more VMs",
	Long:    `By default this command also removes the attached disks, unless the -k flag is passed`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_management.Remove(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(vmCmd, vmLsCmd, vmStartCmd, vmStartAllCmd, vmStopCmd, vmStopAllCmd, vmConsoleCmd, vmRenameCmd, vmRemoveCmd)
	vmCmd.AddCommand(vmLsCmd, vmStartCmd, vmStartAllCmd, vmStopCmd, vmStopAllCmd, vmConsoleCmd, vmRenameCmd, vmRemoveCmd)

	vmConsoleCmd.Flags().BoolVarP(&vm_management.ForceConsoleConnection, "force", "f", false, "Force previous session logout")
	vmRemoveCmd.Flags().BoolVarP(&vm_management.KeepStorage, "keep", "k", false, "Keep VM disk after removal")
}
