// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storagepool_mgt/types.go

package storagepool_mgt

import "vmman4/volume_mgt"

// KeepPoolStorage mirrors vm_mgt.KeepStorage: when true, RemovePool only
// undefines the pool(s) and leaves the underlying storage (directory/disk)
// untouched.
var KeepPoolStorage bool

// StoragePoolInfo describes a single libvirt storage pool.
type StoragePoolInfo struct {
	Name       string
	UUID       string
	State      string // "running", "inactive", …
	TargetPath string // mount point / directory (empty for network pools)
	Volumes    []volume_mgt.VolumeInfo
}
