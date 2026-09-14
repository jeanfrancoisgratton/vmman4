// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/clusterOps.go

package cluster_mgt

import (
	"vmman4/snapshot_mgt"
	"vmman4/vm_mgt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// StartCluster starts every node in cluster, reusing vm_mgt's own start logic.
func StartCluster(cluster string) *ce.CustomError {
	nodes, cerr := nodesFor(cluster)
	if cerr != nil {
		return cerr
	}
	return vm_mgt.StartVM(nodes)
}

// StopCluster stops every node in cluster, reusing vm_mgt's own stop logic.
func StopCluster(cluster string) *ce.CustomError {
	nodes, cerr := nodesFor(cluster)
	if cerr != nil {
		return cerr
	}
	return vm_mgt.StopVM(nodes)
}

// ResetCluster power-cycles every node in cluster (stop then start), reusing
// the same mechanics as `vm reset`.
func ResetCluster(cluster string) *ce.CustomError {
	nodes, cerr := nodesFor(cluster)
	if cerr != nil {
		return cerr
	}
	return vm_mgt.ResetVM(nodes)
}

// SnaprevCluster reverts every node in cluster to snapshotName. An empty
// snapshotName means "the current snapshot", resolved per-node by
// snapshot_mgt.RevertSnapshot.
func SnaprevCluster(cluster, snapshotName string) *ce.CustomError {
	nodes, cerr := nodesFor(cluster)
	if cerr != nil {
		return cerr
	}
	for _, node := range nodes {
		if cerr := snapshot_mgt.RevertSnapshot(node, snapshotName); cerr != nil {
			return cerr
		}
	}
	return nil
}
