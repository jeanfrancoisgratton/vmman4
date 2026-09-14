// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/volumeCreate.go

package volume_mgt

import (
	"fmt"
	"strconv"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// NewVolume is the CLI entry point for `vol create`: it resolves the
// connection (local or remote hypervisor, same as every other *_mgt command),
// parses sizeGB, and creates the volume via CreateVolume. CreateVolume itself
// stays the lower-level primitive, since it's also called internally (e.g. by
// vm_mgt when auto-creating a VM's disks) against an already-open connection.
func NewVolume(poolName, volName, sizeGB string) *ce.CustomError {
	size, e := strconv.ParseFloat(sizeGB, 64)
	if e != nil {
		return &ce.CustomError{Title: "Invalid size_gb value", Message: e.Error()}
	}

	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	vol, cerr := CreateVolume(conn, poolName, volName, size)
	if cerr != nil {
		return cerr
	}
	defer vol.Free()

	fmt.Println(hftx.EnabledSign(volName + hftx.Green(" CREATED")))
	return nil
}

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
