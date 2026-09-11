// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storage_mgt/helpers.go
// Original timestamp: 2026/05/30 13:00:10

package storagemanagement

import (
	"encoding/xml"

	"libvirt.org/go/libvirt"
)

// poolStateString maps a libvirt storage pool state to its display string.
func poolStateString(state libvirt.StoragePoolState) string {
	switch state {
	case libvirt.STORAGE_POOL_RUNNING:
		return "running"
	case libvirt.STORAGE_POOL_INACTIVE:
		return "inactive"
	case libvirt.STORAGE_POOL_BUILDING:
		return "building"
	case libvirt.STORAGE_POOL_DEGRADED:
		return "degraded"
	case libvirt.STORAGE_POOL_INACCESSIBLE:
		return "inaccessible"
	default:
		return "unknown"
	}
}

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
		pState = poolStateString(poolInfo.State)
	}
	if xmlDesc, err := pool.GetXMLDesc(0); err == nil {
		var px poolXML
		if xml.Unmarshal([]byte(xmlDesc), &px) == nil && px.Target.Path != "" {
			pTargetPath = px.Target.Path
		}
	}

	return
}
