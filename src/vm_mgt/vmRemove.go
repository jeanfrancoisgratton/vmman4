// vmman3 : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/vm_mgt/vmRemove.go
// 2022-10-22 12:42:35

package vm_mgt

import (
	"fmt"
	"os"
	"vmman4/connection_mgt"
	storagemanagement "vmman4/storage_mgt"

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
	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, vmname := range args {
		//var host string
		domain, err := shared.GetDomain(conn, vmname)
		if err != nil {
			return err
		}
		defer domain.Free()

		// Shut the VM down, if active
		Wait4Shutdown(domain, vmname)
		fmt.Println(vmname + " now shutdown. Proceeding to removal from inventory.")

		// Storage specs must be gathered while the domain is still defined:
		// UndefineFlags() below removes it from libvirt's inventory, and
		// GetStorageSpecs4VM looks the domain up by name.
		if !KeepStorage {
			storageInfo, e := storagemanagement.GetStorageSpecs4VM(vmname, conn)
			if e != nil {
				return e
			}
			if e := removeStorage(conn, storageInfo.Disks); e != nil {
				return e
			}
		}

		derr := domain.UndefineFlags(lv.DOMAIN_UNDEFINE_SNAPSHOTS_METADATA)
		if derr != nil {
			lverr, ok := derr.(lv.Error)
			if ok {
				fmt.Println(lverr.Message)
				os.Exit(-1)
			}
		}

		fmt.Println(hftx.EnabledSign(vmname + hftx.DimRed(" REMOVED")))
		//fmt.Println(hftx.EnabledSign(fmt.Sprintf("VM %s has been removed.", vmname)))
	}
	return nil
}

// removeStorage(): remove the VM's volumes on the hypervisor the connection_mgt points to.
// This goes through libvirt's storage APIs (rather than a local os.Remove) so that
// removal happens on the actual host owning the connection_mgt, not on the machine running
// vmman4 -- which matters when operating against a remote hypervisor.
func removeStorage(conn *lv.Connect, info []storagemanagement.DiskInfo) *ce.CustomError {
	if len(info) == 0 {
		return nil
	}
	for _, disk := range info {
		if disk.Type != "file" && disk.Type != "block" {
			continue
		}
		vol, err := conn.LookupStorageVolByPath(disk.SourcePath)
		if err != nil {
			return &ce.CustomError{Title: "Error looking up volume " + disk.SourcePath, Message: err.Error()}
		}
		derr := vol.Delete(lv.STORAGE_VOL_DELETE_NORMAL)
		vol.Free()
		if derr != nil {
			return &ce.CustomError{Title: "Error removing " + disk.SourcePath, Message: derr.Error()}
		}
	}
	return nil
}
