// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/cluster_helpers.go

package cluster_mgt

import (
	"strings"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

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

// hypervisorKey returns the raw hypervisor name -- the clusters.json
// top-level key -- implied by the -c/--connectionfile flag, taken as typed
// and stripped of any .json suffix, without running it through
// connection_mgt.ResolveConnectionURI(). Defaults to "local" when -c is unset,
// which lines up with ResolveConnectionURI() leaving shared.ConnectURI at its
// "qemu:///system" default in that same case.
func hypervisorKey() string {
	if shared.ConnectionFilename == "" {
		return DefaultHypervisorKey
	}
	return strings.TrimSuffix(shared.ConnectionFilename, ".json")
}
