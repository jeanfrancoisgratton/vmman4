// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/snpshot_ops.go
// Original timestamp : 2026.11.23 19:13:55

package cluster_mgt

import (
	"vmman4/snapshot_mgt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// SnapshotRevertOnClusterNodes reverts every node in cluster to snapshotName. An empty
// snapshotName means "the current snapshot", resolved per-node by
// snapshot_mgt.RevertSnapshot.
func SnapshotRevertOnClusterNodes(cluster, snapshotName string) *ce.CustomError {
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

// SnapshotRemoveOnClusterNodes removes a snapshot from every node in the cluster. An empty
// snapshotName means "the current snapshot", resolved per-node by
func SnapshotRemoveOnClusterNodes(cluster, snapshotName string) *ce.CustomError {
	nodes, cerr := nodesFor(cluster)
	if cerr != nil {
		return cerr
	}
	for _, node := range nodes {
		if cerr := snapshot_mgt.RemoveSnapshot(node, snapshotName); cerr != nil {
			return cerr
		}
	}
	return nil
}

// SnapshotCreateOnClusterNodes creates a snapshot on every node in the cluster.
func SnapshotCreateOnClusterNodes(cluster, snapshotName string) *ce.CustomError {
	nodes, cerr := nodesFor(cluster)
	if cerr != nil {
		return cerr
	}
	for _, node := range nodes {
		if cerr := snapshot_mgt.CreateSnapshot(node, snapshotName); cerr != nil {
			return cerr
		}
	}
	return nil
}
