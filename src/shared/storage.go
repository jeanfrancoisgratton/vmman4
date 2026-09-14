// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/shared/storage.go

package shared

import (
	"encoding/xml"

	"libvirt.org/go/libvirt"
)

// PoolStateString maps a libvirt storage pool state to its display string.
// Lives here (rather than in storagepool_mgt or volume_mgt) because both
// packages need it and neither may import the other.
func PoolStateString(state libvirt.StoragePoolState) string {
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

type poolXML struct {
	XMLName xml.Name      `xml:"pool"`
	Target  poolTargetXML `xml:"target"`
}

type poolTargetXML struct {
	Path string `xml:"path"`
}

// PoolTargetPath returns a storage pool's target path (mount point /
// directory), or "" for network pools / on error. Shared because both
// storagepool_mgt (listing pools) and volume_mgt (resolving a volume's
// owning pool) need it.
func PoolTargetPath(pool *libvirt.StoragePool) string {
	xmlDesc, err := pool.GetXMLDesc(0)
	if err != nil {
		return ""
	}
	var px poolXML
	if xml.Unmarshal([]byte(xmlDesc), &px) != nil {
		return ""
	}
	return px.Target.Path
}
