// vmman3 : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/vmmanagement/vmRemove.go
// 2022-10-22 12:42:35

package vmmanagement

import (
	"fmt"
	"os"
	storagemanagement "vmman4/storageManagement"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"

	//"github.com/jackc/pgx/v5"
	lv "libvirt.org/go/libvirt"
	//"vmman3/db"
	//"vmman3/env"
	"vmman4/shared"
)

// This will remove the VM, and optionally leave its storage there
func Remove(args []string) *ce.CustomError {
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, vmname := range args {
		//var host string
		domain, err := shared.GetDomain(conn, vmname)
		if domain == nil {
			return nil
		}
		if err != nil {
			return err
		}
		defer domain.Free()

		// Shut the VM down, if active
		Wait4Shutdown(domain, vmname)
		fmt.Println(vmname + " now shutdown. Proceeding to removal from inventory.")
		derr := domain.UndefineFlags(lv.DOMAIN_UNDEFINE_SNAPSHOTS_METADATA)
		if derr != nil {
			lverr, ok := derr.(lv.Error)
			if ok {
				fmt.Println(lverr.Message)
				os.Exit(-1)
			}
		}
		if !KeepStorage {
			storageInfo, e := storagemanagement.GetStorageSpecs4VM(vmname, conn)
			if e != nil {
				return e
			}
			if e := removeStorage(storageInfo.Disks); e != nil {
				return e
			}
		}

		fmt.Println(hftx.EnabledSign(fmt.Sprintf("VM %s has been removed.", vmname)))
	}
	return nil
}

// removeStorage(): remove the VM files from the disks
func removeStorage(info []storagemanagement.DiskInfo) *ce.CustomError {
	if len(info) == 0 {
		return nil
	}
	for _, disk := range info {
		if e := os.Remove(disk.SourcePath); e != nil {
			return &ce.CustomError{Title: "Error removing " + disk.SourcePath, Message: e.Error()}
		}
	}
	return nil
}