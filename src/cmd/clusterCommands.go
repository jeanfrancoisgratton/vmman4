// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/clusterCommands.go

package cmd

import (
	"fmt"

	"vmman4/cluster_mgt"

	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:     "cluster",
	Example: "vmman cluster ls\nvmman cluster snaprev baseline cluster1",
	Short:   "Cluster subcommands",
	Long:    `You need to provide one of the subcommands: ls, define, remove, start, stop, reset, snaprev.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: ls, define, remove, start, stop, reset or snaprev")
	},
}

var clusterLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Example: "vmman cluster ls",
	Short:   "Pretty-print the contents of clusters.json",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cluster_mgt.ListClusters(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var clusterDefineCmd = &cobra.Command{
	Use:     "define <node1> [node2...nodeN] <cluster>",
	Example: "vmman [-c hypervisor] cluster define node1 node2 cluster1",
	Short:   "Define (or redefine) a cluster's node list",
	Long:    `The target hypervisor is taken from -c/--connectionfile (defaults to "local" when omitted). If the hypervisor/cluster pair already exists, its node list is replaced in place.`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		cluster := args[len(args)-1]
		nodes := args[:len(args)-1]
		if err := cluster_mgt.DefineCluster(cluster, nodes); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var clusterRemoveCmd = &cobra.Command{
	Use:     "remove <cluster>",
	Aliases: []string{"rm"},
	Example: "vmman [-c hypervisor] cluster remove cluster1",
	Short:   "Remove a cluster",
	Long:    `The target hypervisor is taken from -c/--connectionfile (defaults to "local" when omitted).`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cluster_mgt.RemoveCluster(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var clusterStartCmd = &cobra.Command{
	Use:     "start <cluster>",
	Aliases: []string{"up"},
	Example: "vmman [-c hypervisor] cluster start cluster_name",
	Short:   "Start every node in a cluster",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cluster_mgt.StartCluster(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var clusterStopCmd = &cobra.Command{
	Use:     "stop <cluster>",
	Aliases: []string{"down"},
	Example: "vmman [-c hypervisor] cluster stop cluster1",
	Short:   "Stop every node in a cluster",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cluster_mgt.StopCluster(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var clusterResetCmd = &cobra.Command{
	Use:     "reset <cluster>",
	Aliases: []string{"reboot", "restart"},
	Example: "vmman [-c hypervisor] cluster reset cluster1",
	Short:   "Stop then start every node in a cluster",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cluster_mgt.ResetCluster(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var clusterSnaprevCmd = &cobra.Command{
	Use:     "snaprev [snapshot_name] <cluster>",
	Aliases: []string{"revert"},
	Example: "vmman [-c hypervisor] cluster snaprev baseline cluster1",
	Short:   "Revert every node in a cluster to a snapshot",
	Long:    "If snapshot_name is not set, each node's current snapshot will be used.",
	Args:    cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		cluster := args[len(args)-1]
		sname := ""
		if len(args) == 2 {
			sname = args[0]
		}
		if err := cluster_mgt.SnaprevCluster(cluster, sname); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterLsCmd, clusterDefineCmd, clusterRemoveCmd, clusterStartCmd, clusterStopCmd, clusterResetCmd, clusterSnaprevCmd)
}
