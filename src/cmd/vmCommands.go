// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/vmCommands.go
// Original timestamp: 2026/05/26 21:21:12

package cmd

import (
	"fmt"

	"vmman4/vm_mgt"

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
		if err := vm_mgt.ListVMs(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show detailed information about a VM",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.VmInfo(args[0]); err != nil {
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
		if err := vm_mgt.StartVM(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStartAllCmd = &cobra.Command{
	Use:   "startall",
	Short: "Start all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.StartAll(); err != nil {
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
		if err := vm_mgt.StopVM(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStopAllCmd = &cobra.Command{
	Use:   "stopall",
	Short: "Stop all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.StopAll(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmResetCmd = &cobra.Command{
	Use:     "reset",
	Aliases: []string{"reboot"},
	Short:   "Stop one or many VMs",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.ResetVM(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmResetAllCmd = &cobra.Command{
	Use:     "resetall",
	Aliases: []string{"rebootall"},
	Short:   "Stop/Start all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.ResetAllVMs(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmConsoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Open a console session on the VM",
	Long:  "a -f flag will force the connection_mgt to the console, if that connection_mgt was already opened.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.Console(args[0]); err != nil {
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
		if err := vm_mgt.RenameVM(args[0], args[1]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmSetMemCmd = &cobra.Command{
	Use:   "setmem",
	Short: "Set a VM's memory (in MiB)",
	Long: `Expects 1 or 2 numeric arguments: min_mem [max_mem].
If only min_mem is passed, min_mem = max_mem.
Overcommitting the hypervisor's physical memory is warned about, but not blocked.`,
	Args: cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.SetVMem(args[0], args[1:]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmSetVcpusCmd = &cobra.Command{
	Use:     "setvcpus",
	Aliases: []string{"setcpu", "setcpus"},
	Short:   "Set a VM's vCPU count",
	Long:    `Overcommitting the hypervisor's physical CPUs is warned about, but not blocked.`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.SetVCPUs(args[0], args[1]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmDumpXmlCmd = &cobra.Command{
	Use:   "dumpxml",
	Short: "Dump a VM's XML configuration to a file",
	Long:  `Shuts the VM down (if active) before dumping its inactive, migratable XML description.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.DumpVmXML(args[0], args[1]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmCreateCmd = &cobra.Command{
	Use:   "create <spec.json>",
	Short: "Create (define) a VM from a JSON spec file",
	Long: `Reads a VMSpec JSON file and defines the domain on the target hypervisor. The VM is not started; use 'vm start' to boot it.

Pass -s/--sample [outfile] instead of a real spec file to generate a fully annotated example spec (defaults to ~/.config/JFG/vmman4/vmspec.sample.json) -- useful if you're not familiar with libvirt's concepts (pools, volumes, networks, etc).`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if vm_mgt.SampleMode {
			dest := vm_mgt.DefaultSampleFile()
			if len(args) == 1 {
				dest = args[0]
			}
			if err := vm_mgt.WriteSample(dest); err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		if len(args) != 1 {
			fmt.Println("Usage: vm create <spec.json>  (or: vm create -s [outfile] to generate a sample)")
			return
		}
		if err := vm_mgt.CreateVM(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmProvisionCmd = &cobra.Command{
	Use:   "provision TARGET_HOSTNAME IP_ADDRESS TEMPLATE_NAME",
	Short: "Clone a template VM into a new, network-provisioned VM",
	Long: `Clones TEMPLATE_NAME's disk (TEMPLATE_NAME.qcow2) and domain definition into a
new VM named TARGET_HOSTNAME, boots it, then uses the QEMU guest agent
(must already be installed and running in the template) to set its hostname
and network configuration and regenerate its SSH host keys.

Network defaults (storage pool, CIDR, gateway, DNS) come from an environment
JSON file -- see ~/.config/JFG/vmman4/env-sample.json for a baseline; pass
-E to use a file other than the default ~/.config/JFG/vmman4/env.json.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.ProvisionVM(args[0], args[1], args[2], vm_mgt.EnvFile); err != nil {
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
		if err := vm_mgt.RemoveVM(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(vmCmd, vmLsCmd, vmStartCmd, vmStopCmd)
	vmCmd.AddCommand(vmLsCmd, vmInfoCmd, vmCreateCmd, vmProvisionCmd, vmStartCmd, vmStartAllCmd, vmStopCmd, vmStopAllCmd, vmResetCmd,
		vmResetAllCmd, vmConsoleCmd, vmRenameCmd, vmRemoveCmd, vmSetMemCmd, vmSetVcpusCmd, vmDumpXmlCmd)

	vmConsoleCmd.Flags().BoolVarP(&vm_mgt.ForceConsoleConnection, "force", "f", false, "Force previous session logout")
	vmRemoveCmd.Flags().BoolVarP(&vm_mgt.KeepStorage, "keep", "k", false, "Keep VM disk after removal")
	vmCreateCmd.Flags().BoolVarP(&vm_mgt.SampleMode, "sample", "s", false, "Write an annotated sample spec file instead of creating a VM")
	vmProvisionCmd.Flags().StringVarP(&vm_mgt.EnvFile, "environment", "E", "", "Environment file (default: ~/.config/JFG/vmman4/env.json)")
}
