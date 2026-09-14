// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/volumeRemove.go

package volume_mgt

import (
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// RemoveVolume deletes one or more volumes from the named storage pool.
func RemoveVolume(poolName string, volNames []string) *ce.CustomError {
	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	pool, perr := conn.LookupStoragePoolByName(poolName)
	if perr != nil {
		return &ce.CustomError{Title: fmt.Sprintf("RemoveVolume: pool %q not found", poolName), Message: perr.Error()}
	}
	defer pool.Free()

	for _, volName := range volNames {
		vol, verr := pool.LookupStorageVolByName(volName)
		if verr != nil {
			return &ce.CustomError{Title: fmt.Sprintf("RemoveVolume: volume %q not found in pool %q", volName, poolName), Message: verr.Error()}
		}

		derr := vol.Delete(libvirt.STORAGE_VOL_DELETE_NORMAL)
		vol.Free()
		if derr != nil {
			return &ce.CustomError{Title: fmt.Sprintf("RemoveVolume: cannot remove volume %q", volName), Message: derr.Error()}
		}

		fmt.Println(hftx.EnabledSign(volName + hftx.DimRed(" REMOVED")))
	}
	return nil
}
