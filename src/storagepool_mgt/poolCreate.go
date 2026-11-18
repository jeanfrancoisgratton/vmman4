// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storagepool_mgt/poolCreate.go

package storagepool_mgt

import (
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// CreatePool defines a new directory-backed ("dir" type) storage pool at
// targetPath, marks it to autostart, and builds + starts it immediately.
func CreatePool(poolName, targetPath string) *ce.CustomError {
	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	poolXML := fmt.Sprintf(`<pool type="dir">
  <name>%s</name>
  <target>
    <path>%s</path>
  </target>
</pool>`, poolName, targetPath)

	pool, derr := conn.StoragePoolDefineXML(poolXML, 0)
	if derr != nil {
		return &ce.CustomError{Title: "CreatePool: cannot define pool " + poolName, Message: derr.Error()}
	}
	defer pool.Free()

	if aerr := pool.SetAutostart(true); aerr != nil {
		return &ce.CustomError{Title: "CreatePool: cannot set autostart on pool " + poolName, Message: aerr.Error()}
	}

	if cerr := pool.Create(libvirt.STORAGE_POOL_CREATE_WITH_BUILD); cerr != nil {
		return &ce.CustomError{Title: "CreatePool: cannot build/start pool " + poolName, Message: cerr.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Pool " + hftx.Brown(poolName) + hftx.Green(" created")))
	}
	return nil
}
