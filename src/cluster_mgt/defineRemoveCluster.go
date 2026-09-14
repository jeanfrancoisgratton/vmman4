// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/defineRemoveCluster.go

package cluster_mgt

import (
	"fmt"

	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// DefineCluster creates (or, if hypervisor/cluster already exists, in-place
// replaces) the node list for a cluster on the hypervisor implied by
// hypervisorKey().
func DefineCluster(cluster string, nodes []string) *ce.CustomError {
	seen := make(map[string]bool, len(nodes))
	for _, node := range nodes {
		if seen[node] {
			return &ce.CustomError{Title: "DefineCluster: duplicate node name", Message: node}
		}
		seen[node] = true
	}

	hv := hypervisorKey()

	cf, cerr := loadClustersFile()
	if cerr != nil {
		return cerr
	}

	if cf[hv] == nil {
		cf[hv] = map[string][]string{}
	}
	cf[hv][cluster] = nodes

	if cerr := saveClustersFile(cf); cerr != nil {
		return cerr
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("%s/%s", hv, cluster) + hftx.Green(" DEFINED") + fmt.Sprintf(" (%d nodes)", len(nodes))))
	}
	return nil
}

// RemoveCluster deletes a cluster from the hypervisor implied by
// hypervisorKey(). The hypervisor key itself is left in place (even if it
// ends up with no clusters left under it), since an empty object is still
// valid JSON and keeps the hypervisor "known".
func RemoveCluster(cluster string) *ce.CustomError {
	hv := hypervisorKey()

	cf, cerr := loadClustersFile()
	if cerr != nil {
		return cerr
	}

	if _, ok := cf[hv]; !ok {
		return &ce.CustomError{Title: "RemoveCluster: unknown hypervisor", Message: hv}
	}
	if _, ok := cf[hv][cluster]; !ok {
		return &ce.CustomError{Title: "RemoveCluster: unknown cluster", Message: cluster + " on hypervisor " + hv}
	}

	delete(cf[hv], cluster)

	if cerr := saveClustersFile(cf); cerr != nil {
		return cerr
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign(fmt.Sprintf("%s/%s", hv, cluster) + hftx.DimRed(" REMOVED")))
	}
	return nil
}
