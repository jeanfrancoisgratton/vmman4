// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/hypervisorKey.go

package cluster_mgt

import (
	"strings"

	"vmman4/shared"
)

// DefaultHypervisorKey is the clusters.json hypervisor key used when
// -c/--connectionfile is not given.
const DefaultHypervisorKey = "local"

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
