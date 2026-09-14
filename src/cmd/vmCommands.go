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
	Use:     "vm",
	Short:   "Virtual machines management",
	Long:    `Commands to manage the VMs`,
	Example: "vmman vm list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.ScrollSign("Available commands: vm [ list | stop | start | restart ]"))
	},
}

var vmLsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "vmman vm list",
	Short:   "List all VMs",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.ListVMs(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmInfoCmd = &cobra.Command{
	Use:     "info <VM>",
	Example: "vmman vm info myvm",
	Short:   "Show detailed information about a VM",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.VmInfo(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStartCmd = &cobra.Command{
	Use:     "start <VM...>",
	Aliases: []string{"up"},
	Example: "vmman vm start myvm\nvmman vm start myvm1 myvm2",
	Short:   "Start one or many VMs",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.StartVM(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStartAllCmd = &cobra.Command{
	Use:     "startall",
	Example: "vmman vm startall",
	Short:   "Start all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.StartAll(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStopCmd = &cobra.Command{
	Use:     "stop <VM...>",
	Aliases: []string{"down"},
	Example: "vmman vm stop myvm\nvmman vm stop myvm1 myvm2",
	Short:   "Stop one or many VMs",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.StopVM(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmStopAllCmd = &cobra.Command{
	Use:     "stopall",
	Example: "vmman vm stopall",
	Short:   "Stop all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.StopAll(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmResetCmd = &cobra.Command{
	Use:     "reset <VM...>",
	Aliases: []string{"reboot"},
	Example: "vmman vm reset myvm",
	Short:   "Stop then start one or many VMs",
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
	Example: "vmman vm resetall",
	Short:   "Stop/Start all VMs at once",
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.ResetAllVMs(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmConsoleCmd = &cobra.Command{
	Use:     "console <VM>",
	Example: "vmman vm console myvm\nvmman vm console myvm -f",
	Short:   "Open a console session on the VM",
	Long:    "The -f/--force flag disconnects a previous session on the same VM before opening the new one.",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.Console(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmRenameCmd = &cobra.Command{
	Use:     "rename <OLD> <NEW>",
	Example: "vmman vm rename myvm newname",
	Short:   "Rename a VM",
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
	Use:     "setmem <VM> <MIN_MB> [MAX_MB]",
	Example: "vmman vm setmem myvm 2048\nvmman vm setmem myvm 2048 4096",
	Short:   "Set a VM's memory (in MiB)",
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
	Use:     "setvcpus <VM> <COUNT>",
	Aliases: []string{"setcpu", "setcpus"},
	Example: "vmman vm setvcpus myvm 4",
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
	Use:     "dumpxml <VM> <FILE>",
	Example: "vmman vm dumpxml myvm myvm.xml",
	Short:   "Dump a VM's XML configuration to a file",
	Long:    `Shuts the VM down (if active) before dumping its inactive, migratable XML description.`,
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := vm_mgt.DumpVmXML(args[0], args[1]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var vmCreateCmd = &cobra.Command{
	Use:     "create <spec.json>",
	Example: "vmman vm create myvm.json\nvmman vm create -s myvm-sample.json",
	Short:   "Create (define) a VM from a JSON spec file",
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
	Use:     "provision TARGET_HOSTNAME IP_ADDRESS TEMPLATE_NAME",
	Example: "vmman vm provision newhost 192.168.1.50 mytemplate",
	Short:   "Clone a template VM into a new, network-provisioned VM",
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
	Use:     "rm <VM...>",
	Aliases: []string{"remove", "destroy", "delete"},
	Example: "vmman vm rm myvm\nvmman vm rm myvm1 myvm2 -k",
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
