// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/helpers.go

package volume_mgt

import (
	"vmman4/shared"

	"libvirt.org/go/libvirt"
)

// resolveVolumeInfo tries to obtain a volume's logical capacity along with
// its owning storage pool's name, target path, and state. It first attempts
// a direct path lookup, then falls back to pool+volume name. Pool fields are
// "n/a" when the disk has no associated storage pool (e.g. network disks).
func resolveVolumeInfo(conn *libvirt.Connect, path, poolName, volName string) (sizeBytes uint64, pName, pTargetPath, pState string) {
	pName, pTargetPath, pState = "n/a", "n/a", "n/a"

	var vol *libvirt.StorageVol
	if path != "" {
		vol, _ = conn.LookupStorageVolByPath(path)
	}
	if vol == nil && poolName != "" && volName != "" {
		if pool, err := conn.LookupStoragePoolByName(poolName); err == nil {
			defer pool.Free()
			vol, _ = pool.LookupStorageVolByName(volName)
		}
	}
	if vol == nil {
		return
	}
	defer vol.Free()

	if info, err := vol.GetInfo(); err == nil {
		sizeBytes = info.Capacity
	}

	pool, err := vol.LookupPoolByVolume()
	if err != nil {
		return
	}
	defer pool.Free()

	if name, err := pool.GetName(); err == nil {
		pName = name
	}
	if poolInfo, err := pool.GetInfo(); err == nil {
		pState = shared.PoolStateString(poolInfo.State)
	}
	if tp := shared.PoolTargetPath(pool); tp != "" {
		pTargetPath = tp
	}

	return
}
