// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/types.go

package cluster_mgt

// ClustersFile maps hypervisor name -> cluster name -> node names.
// A cluster name may appear under several hypervisors independently, since a
// single cluster (eg. a k8s cluster) can be spread across them.
type ClustersFile map[string]map[string][]string
