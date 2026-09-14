// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/volumeCommands.go

package cmd

import (
	"fmt"

	"vmman4/volume_mgt"

	"github.com/spf13/cobra"
)

// volCmd covers all the storage volume-related subcommands
var volCmd = &cobra.Command{
	Use:     "vol",
	Aliases: []string{"volume"},
	Example: "vmman vol list",
	Short:   "Storage volume subcommands",
	Long:    `You need to provide one of the subcommands: ls, create, rm.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, create or rm")
	},
}

var volLsCmd = &cobra.Command{
	Use:     "list [pool]",
	Aliases: []string{"ls"},
	Example: "vmman vol list\nvmman vol list default",
	Short:   "List storage volumes, optionally restricted to one pool",
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pool := ""
		if len(args) == 1 {
			pool = args[0]
		}
		if err := volume_mgt.ListVolumes(pool); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var volCreateCmd = &cobra.Command{
	Use:     "create <pool> <name> <size_gb>",
	Example: "vmman vol create default my-vm.qcow2 20",
	Short:   "Create a new qcow2 volume in a storage pool",
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := volume_mgt.NewVolume(args[0], args[1], args[2]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var volRemoveCmd = &cobra.Command{
	Use:     "rm <pool> <name...>",
	Aliases: []string{"remove", "destroy", "delete"},
	Example: "vmman vol rm default my-vm.qcow2",
	Short:   "Remove one or more volumes from a storage pool",
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := volume_mgt.RemoveVolume(args[0], args[1:]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(volCmd)
	volCmd.AddCommand(volLsCmd, volCreateCmd, volRemoveCmd)
}
