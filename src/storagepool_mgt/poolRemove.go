// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storagepool_mgt/poolRemove.go

package storagepool_mgt

import (
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// KeepPoolStorage mirrors vm_mgt.KeepStorage: when true, RemovePool only
// undefines the pool(s) and leaves the underlying storage (directory/disk)
// untouched.
var KeepPoolStorage bool

// RemovePool stops (if active), optionally wipes the underlying storage of,
// and undefines one or more storage pools. By default the pool's storage is
// wiped, unless KeepPoolStorage is set -- in which case only the pool
// definition is removed.
func RemovePool(poolNames []string) *ce.CustomError {
	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, poolName := range poolNames {
		pool, perr := conn.LookupStoragePoolByName(poolName)
		if perr != nil {
			return &ce.CustomError{Title: fmt.Sprintf("RemovePool: pool %q not found", poolName), Message: perr.Error()}
		}

		isActive, _ := pool.IsActive()

		if !KeepPoolStorage {
			if !isActive {
				pool.Free()
				return &ce.CustomError{
					Title:   fmt.Sprintf("RemovePool: pool %q is not active", poolName),
					Message: "the pool must be active to wipe its storage -- start it first, or pass -k/--keep to only remove the pool definition",
				}
			}
			if derr := pool.Delete(libvirt.STORAGE_POOL_DELETE_NORMAL); derr != nil {
				pool.Free()
				return &ce.CustomError{Title: fmt.Sprintf("RemovePool: cannot delete storage for pool %q", poolName), Message: derr.Error()}
			}
		}

		if isActive {
			if derr := pool.Destroy(); derr != nil {
				pool.Free()
				return &ce.CustomError{Title: fmt.Sprintf("RemovePool: cannot stop pool %q", poolName), Message: derr.Error()}
			}
		}

		uerr := pool.Undefine()
		pool.Free()
		if uerr != nil {
			return &ce.CustomError{Title: fmt.Sprintf("RemovePool: cannot undefine pool %q", poolName), Message: uerr.Error()}
		}

		fmt.Println(hftx.EnabledSign(poolName + hftx.DimRed(" REMOVED")))
	}
	return nil
}
