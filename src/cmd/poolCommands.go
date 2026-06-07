// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/poolCommands.go
// Original timestamp: 2026/05/30 21:52:04

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	storagemanagement "vmman4/storageManagement"
)

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Storage pool subcommands",
	Long:  `You need to provide one of the subcommands: ls, create, rm, info.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, create, rm or info")
	},
}

var poolLsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "vmman pool list",
	Short:   "Gives extended information about the storage pools",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := storagemanagement.ListStoragePools(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

//var storageLsCmd = &cobra.Command{
//	Use:     "slist",
//	Aliases: []string{"sls"},
//	Example: "vmman slist VIRTUAL_MACHINE",
//	Short:   "Gives extended information about the storage of a specific VM",
//	Args:    cobra.MinimumNArgs(1),
//	Run: func(cmd *cobra.Command, args []string) {
//		if err := storagemanagement.ListVMStorage(args[0], true); err != nil {
//			fmt.Println(err.Error())
//		}
//	},
//}

func init() {
	rootCmd.AddCommand(poolCmd /*storageLsCmd*/)

	poolCmd.AddCommand(poolLsCmd)
}
