// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/snapshotCommands.go
// Original timestamp : 2026.09.12 18:58:48

package cmd

import (
	"fmt"
	"vmman4/snapshot_mgt"

	"github.com/spf13/cobra"
)

var snapCreateDescription, snapDumpXMLFile string
var snapLsAsTree, snapRmWithChildren, snapRmOnlyChildren bool

var snapCmd = &cobra.Command{
	Use:     "snapshot",
	Aliases: []string{"snap"},
	Short:   "Snapshot subcommands",
	Long:    `You need to provide one of the subcommands: ls, create, rm, revert`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, create, rm or dumpxml")
	},
}

var snapLsCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Example: "vmman snap list VM",
	Short:   "Lists all snapshots for one or more VMs",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := snapshot_mgt.ListSnapshots(args, snapLsAsTree); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var snapCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Example: "vmman snap create VM SNAPSHOTNAME [PARENTSNAP]",
	Short:   "Creates a snapshot SNAPSHOTNAME off the PARENTSNAP snapshot on VM",
	Long:    "If PARENTSNAP is not set, the current snapshot will be used",
	Args:    cobra.MaximumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		pname := ""
		if len(args) > 2 {
			pname = args[2]
		}
		if err := snapshot_mgt.CreateSnapshot(args[0], args[1], pname, snapCreateDescription); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var snapRevertCmd = &cobra.Command{
	Use:     "revert",
	Aliases: []string{"set"},
	Example: "vmman snap revert VM [SNAPSHOTNAME]",
	Short:   "Reverts (sets) a VM to SNAPSHOTNAME",
	Long:    "If SNAPSHOTNAME is not set, the current snapshot will be used",
	Args:    cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sname := ""
		if len(args) > 1 {
			sname = args[1]
		}
		if err := snapshot_mgt.RevertSnapshot(args[0], sname); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var snapRmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove"},
	Example: "vmman snap rm VM [SNAPSHOTNAME]",
	Short:   "Removes a snapshot (SNAPSHOTNAME) from a VM",
	Long:    "If SNAPSHOTNAME is not set, the current snapshot will be used",
	Args:    cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sname := ""
		if len(args) > 1 {
			sname = args[1]
		}
		if err := snapshot_mgt.RemoveSnapshot(args[0], sname, snapRmWithChildren, snapRmOnlyChildren); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var snapDumpXmlCmd = &cobra.Command{
	Use:     "dumpxml",
	Example: "vmman snap dumpxml VM [SNAPSHOTNAME] [-f filename]",
	Short:   "Dumps a snapshot's XML description",
	Long:    "If SNAPSHOTNAME is not set, the current snapshot will be used. Without -f, the XML is printed to stdout.",
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		sname := ""
		if len(args) > 1 {
			sname = args[1]
		}
		if err := snapshot_mgt.DumpSnapshotXML(args[0], sname, snapDumpXMLFile); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(snapCmd)
	snapCmd.AddCommand(snapLsCmd, snapCreateCmd, snapRevertCmd, snapRmCmd, snapDumpXmlCmd)

	snapLsCmd.Flags().BoolVarP(&snapLsAsTree, "tree", "x", false, "Render the snapshot hierarchy as a tree")

	snapCreateCmd.Flags().StringVarP(&snapCreateDescription, "description", "d", "", "Description to embed in the snapshot")

	snapRmCmd.Flags().BoolVarP(&snapRmWithChildren, "with-children", "k", false, "Also delete the snapshot's children")
	snapRmCmd.Flags().BoolVarP(&snapRmOnlyChildren, "only-children", "K", false, "Delete only the snapshot's children, keeping the snapshot itself")
	snapRmCmd.MarkFlagsMutuallyExclusive("with-children", "only-children")

	snapDumpXmlCmd.Flags().StringVarP(&snapDumpXMLFile, "file", "f", "", "Write XML to file instead of stdout")
}
