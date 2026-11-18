// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storagepool_mgt/remove.go

package storagepool_mgt

import (
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// RemovePool stops (if active), optionally wipes the underlying storage of,
// and undefines one or more storage pools. By default, the pool's storage is
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
			return &ce.CustomError{Title: "RemovePool: pool " + poolName + "not found", Message: perr.Error()}
		}

		isActive, _ := pool.IsActive()

		if !KeepPoolStorage {
			if !isActive {
				pool.Free()
				return &ce.CustomError{Title: "RemovePool: pool " + poolName + " is not active",
					Message: "the pool must be active to wipe its storage -- start it first, or pass -k/--keep to only remove the pool definition",
				}
			}
			if derr := pool.Delete(libvirt.STORAGE_POOL_DELETE_NORMAL); derr != nil {
				pool.Free()
				return &ce.CustomError{Title: "RemovePool: cannot delete storage for pool " + poolName, Message: derr.Error()}
			}
		}

		if isActive {
			if derr := pool.Destroy(); derr != nil {
				pool.Free()
				return &ce.CustomError{Title: "RemovePool: cannot stop pool " + poolName, Message: derr.Error()}
			}
		}

		uerr := pool.Undefine()
		pool.Free()
		if uerr != nil {
			return &ce.CustomError{Title: "RemovePool: cannot undefine pool " + poolName, Message: uerr.Error()}
		}

		if !shared.QuietOutput {
			fmt.Println(hftx.EnabledSign("Pool " + hftx.Brown(poolName) + hftx.DimRed(" removed")))
		}
	}
	return nil
}
