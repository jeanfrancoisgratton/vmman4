// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/clusterCommands.go

package cmd

import (
	"fmt"

	"vmman4/cluster_mgt"

	"github.com/spf13/cobra"
)

// clusterCmd's own Run doubles as the dispatcher for `vmman [-c hv] CLUSTER_NAME snaprev [SNAPSHOT_NAME]`:
// cobra only routes to a registered child (ls/define/remove/start/stop/reset) when args[0] matches its
// Use/Alias, so anything else falls through here with args[0] being the cluster name.
// NOTE: a cluster literally named "ls", "define", "remove", "start", "up", "stop", "down", "reset" or
// "reboot" can never reach this path -- cobra will always treat it as the matching subcommand instead.
var clusterCmd = &cobra.Command{
	Use:     "cluster",
	Example: "vmman cluster ls\nvmman cluster cluster1 snaprev baseline",
	Short:   "Cluster subcommands",
	Long: `You need to provide one of the subcommands: ls, define, remove, start, stop, reset.
To revert a cluster's nodes to a snapshot: vmman cluster CLUSTER_NAME snaprev [SNAPSHOT_NAME]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 2 && (args[1] == "snaprev" || args[1] == "revert") {
			sname := ""
			if len(args) > 2 {
				sname = args[2]
			}
			if err := cluster_mgt.SnaprevCluster(args[0], sname); err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println("You need to provide one of the following subcommands: ls, define, remove, start, stop or reset")
		fmt.Println(`To revert a cluster's nodes to a snapshot: vmman cluster CLUSTER_NAME snaprev [SNAPSHOT_NAME]`)
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
	Use:     "define <cluster> <node1> [node2...nodeN]",
	Example: "vmman [-c hypervisor] cluster define cluster1 node1 node2",
	Short:   "Define (or redefine) a cluster's node list",
	Long:    `The target hypervisor is taken from -c/--connectionfile (defaults to "local" when omitted). If the hypervisor/cluster pair already exists, its node list is replaced in place.`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cluster_mgt.DefineCluster(args[0], args[1:]); err != nil {
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
	Example: "vmman [-c hypervisor] cluster start cluster1",
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
	Aliases: []string{"reboot"},
	Example: "vmman [-c hypervisor] cluster reset cluster1",
	Short:   "Stop then start every node in a cluster",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cluster_mgt.ResetCluster(args[0]); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterLsCmd, clusterDefineCmd, clusterRemoveCmd, clusterStartCmd, clusterStopCmd, clusterResetCmd)
}
