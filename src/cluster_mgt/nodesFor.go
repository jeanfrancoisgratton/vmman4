// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/nodesFor.go

package cluster_mgt

import ce "github.com/jeanfrancoisgratton/customError/v3"

// nodesFor returns the node list for cluster on the hypervisor implied by
// hypervisorKey().
func nodesFor(cluster string) ([]string, *ce.CustomError) {
	hv := hypervisorKey()

	cf, cerr := loadClustersFile()
	if cerr != nil {
		return nil, cerr
	}

	clusters, ok := cf[hv]
	if !ok {
		return nil, &ce.CustomError{Title: "Unknown hypervisor", Message: hv}
	}
	nodes, ok := clusters[cluster]
	if !ok {
		return nil, &ce.CustomError{Title: "Unknown cluster", Message: cluster + " on hypervisor " + hv}
	}
	return nodes, nil
}
