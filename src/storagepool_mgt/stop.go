// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storagepool_mgt/stop.go

package storagepool_mgt

import (
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// StopPool deactivates (stops) an active storage pool. This does not remove
// its definition or delete any volumes; use `pool rm` (not yet implemented)
// to undefine a pool entirely.
func StopPool(poolName string) *ce.CustomError {
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
		return &ce.CustomError{Title: "StopPool: pool " + poolName + " not found", Message: perr.Error()}
	}
	defer pool.Free()

	if derr := pool.Destroy(); derr != nil {
		return &ce.CustomError{Title: "StopPool: cannot stop pool " + poolName, Message: derr.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Pool " + hftx.Brown(poolName) + hftx.DimRed(" stopped")))
	}
	return nil
}
