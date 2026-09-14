// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storagepool_mgt/poolStart.go

package storagepool_mgt

import (
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// StartPool activates (starts) an already-defined but inactive storage pool.
func StartPool(poolName string) *ce.CustomError {
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
		return &ce.CustomError{Title: fmt.Sprintf("StartPool: pool %q not found", poolName), Message: perr.Error()}
	}
	defer pool.Free()

	if cerr := pool.Create(0); cerr != nil {
		return &ce.CustomError{Title: fmt.Sprintf("StartPool: cannot start pool %q", poolName), Message: cerr.Error()}
	}

	fmt.Println(hftx.EnabledSign(poolName + hftx.Green(" STARTED")))
	return nil
}
