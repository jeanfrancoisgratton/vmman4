// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storageManagement/helpers.go
// Original timestamp: 2026/05/30 13:00:10

package storagemanagement

import (
	"libvirt.org/go/libvirt"
)

// resolveVolumeSize tries to obtain a volume's logical capacity.
// It first attempts a direct path lookup, then falls back to pool+volume name.
func resolveVolumeSize(conn *libvirt.Connect, path, poolName, volName string) uint64 {
	if path != "" {
		if vol, err := conn.LookupStorageVolByPath(path); err == nil {
			defer vol.Free()
			if info, err := vol.GetInfo(); err == nil {
				return info.Capacity
			}
		}
	}
	if poolName != "" && volName != "" {
		pool, err := conn.LookupStoragePoolByName(poolName)
		if err != nil {
			return 0
		}
		defer pool.Free()
		vol, err := pool.LookupStorageVolByName(volName)
		if err != nil {
			return 0
		}
		defer vol.Free()
		if info, err := vol.GetInfo(); err == nil {
			return info.Capacity
		}
	}
	return 0
}
