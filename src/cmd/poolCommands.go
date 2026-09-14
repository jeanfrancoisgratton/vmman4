// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/poolCommands.go
// Original timestamp: 2026/05/30 21:52:04

package cmd

import (
	"fmt"

	"vmman4/storagepool_mgt"

	"github.com/spf13/cobra"
)

var poolCmd = &cobra.Command{
	Use:     "pool",
	Example: "vmman pool list",
	Short:   "Storage pool subcommands",
	Long:    `You need to provide one of the subcommands: ls, create, start, stop, rm, info.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, create, start, stop, rm or info")
	},
}

var poolLsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "vmman pool list",
	Short:   "Gives extended information about the storage pools",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := storagepool_mgt.ListStoragePools(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var poolCreateCmd = &cobra.Command{
	Use:     "create <name> <target_path>",
	Example: "vmman pool create extra /data/vmpool",
	Short:   "Define, build and start a new directory-backed storage pool",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := storagepool_mgt.CreatePool(args[0], args[1]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var poolStartCmd = &cobra.Command{
	Use:     "start <name>",
	Aliases: []string{"up"},
	Example: "vmman pool start extra",
	Short:   "Start (activate) a defined but inactive storage pool",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := storagepool_mgt.StartPool(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var poolStopCmd = &cobra.Command{
	Use:     "stop <name>",
	Aliases: []string{"down"},
	Example: "vmman pool stop extra",
	Short:   "Stop (deactivate) an active storage pool",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := storagepool_mgt.StopPool(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var poolRemoveCmd = &cobra.Command{
	Use:     "rm <name...>",
	Aliases: []string{"remove", "destroy", "delete"},
	Example: "vmman pool rm extra",
	Short:   "Stop, wipe and undefine one or more storage pools",
	Long:    `By default this command also wipes the pool's underlying storage, unless the -k flag is passed`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := storagepool_mgt.RemovePool(args); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(poolCmd)

	poolCmd.AddCommand(poolLsCmd, poolCreateCmd, poolStartCmd, poolStopCmd, poolRemoveCmd)

	poolRemoveCmd.Flags().BoolVarP(&storagepool_mgt.KeepPoolStorage, "keep", "k", false, "Keep pool storage after removal")
}
