// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/volumeCreate.go

package volume_mgt

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// CreateVolume creates a new qcow2 volume of the given capacity (in GiB)
// inside the named storage pool. The caller must Free() the returned volume.
func CreateVolume(conn *libvirt.Connect, poolName, volName string, sizeGB float64) (*libvirt.StorageVol, *ce.CustomError) {
	if sizeGB <= 0 {
		return nil, &ce.CustomError{
			Title:   fmt.Sprintf("CreateVolume: invalid size for volume %q", volName),
			Message: fmt.Sprintf("size_gb must be > 0, got %v", sizeGB),
		}
	}

	pool, err := conn.LookupStoragePoolByName(poolName)
	if err != nil {
		return nil, &ce.CustomError{Title: fmt.Sprintf("CreateVolume: pool %q not found", poolName), Message: err.Error()}
	}
	defer pool.Free()

	sizeBytes := uint64(sizeGB * 1024 * 1024 * 1024)

	volXML := fmt.Sprintf(`<volume>
  <name>%s</name>
  <capacity unit="bytes">%d</capacity>
  <target>
    <format type="qcow2"/>
  </target>
</volume>`, volName, sizeBytes)

	vol, err := pool.StorageVolCreateXML(volXML, 0)
	if err != nil {
		return nil, &ce.CustomError{
			Title:   fmt.Sprintf("CreateVolume: cannot create volume %q in pool %q", volName, poolName),
			Message: err.Error(),
		}
	}

	return vol, nil
}
