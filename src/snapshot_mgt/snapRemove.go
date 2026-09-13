// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/snapshot_mgt/snapRemove.go
// Original timestamp : 2026.09.13 12:57:50

package snapshot_mgt

import (
	"errors"
	"fmt"
	"strings"
	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

func RemoveSnapshot(vmname, sname string, withChildren, onlyChildren bool) *ce.CustomError {
	var cerr *ce.CustomError
	var conn *libvirt.Connect
	var domain *libvirt.Domain

	if withChildren && onlyChildren {
		return &ce.CustomError{Title: "Conflicting flags", Message: "--with-children and --only-children are mutually exclusive"}
	}

	if cerr = connection_mgt.ResolveConnectionURI(); cerr != nil {
		return cerr
	}
	if conn, cerr = shared.Connect2HVM(); cerr != nil {
		return cerr
	}
	defer conn.Close()

	if domain, cerr = shared.GetDomain(conn, vmname); cerr != nil {
		return cerr
	}
	defer domain.Free()
	numsnap, _ := domain.SnapshotNum(0)
	// no snapshot, nothing to do
	if numsnap == 0 {
		return &ce.CustomError{Fatality: ce.NotAnError, Message: "Snapshot number is zero"}
	}
	// no snapshot name were provided, we'll use the current snapshot
	if sname == "" {
		if sname, cerr = GetCurrentSnapshotName(conn, vmname); cerr != nil {
			return cerr
		}
	}
	// there are snapshots, so we need to shut the vm down before proceeding
	shared.Wait4Shutdown(domain, vmname)

	snapshot, err := domain.SnapshotLookupByName(sname, 0)
	if err != nil {
		var lverr libvirt.Error
		if errors.As(err, &lverr) {
			if strings.HasPrefix(lverr.Message, "Domain snapshot not found") {
				return &ce.CustomError{Title: "Snapshot not found",
					Message: "No snapshot named " + hftx.Bold(hftx.Red(sname)) + " exists for " + hftx.Bold(hftx.Red(vmname))}
			}
			return &ce.CustomError{Title: "snapshot_mgt.RevertSnapshot() failure", Message: lverr.Message}
		}
		return &ce.CustomError{Title: "snapshot_mgt.RevertSnapshot() failure", Message: err.Error()}
	}
	defer snapshot.Free()

	var flags libvirt.DomainSnapshotDeleteFlags
	switch {
	case withChildren:
		flags = libvirt.DOMAIN_SNAPSHOT_DELETE_CHILDREN
	case onlyChildren:
		flags = libvirt.DOMAIN_SNAPSHOT_DELETE_CHILDREN_ONLY
	}

	if err = snapshot.Delete(flags); err != nil {
		var lverr libvirt.Error
		if errors.As(err, &lverr) {
			return &ce.CustomError{Title: "Cannot delete the " + sname + " snapshot from " + vmname, Message: lverr.Message}
		}
		return &ce.CustomError{Title: "Cannot delete the " + sname + " snapshot from " + vmname, Message: err.Error()}
	}

	verb := "removed"
	if onlyChildren {
		verb = "pruned of its children"
	} else if withChildren {
		verb = "removed along with its children"
	}
	fmt.Println(hftx.EnabledSign(fmt.Sprintf("Snapshot %s %s from %s", sname,
		hftx.Bold(hftx.DimRed(verb)), vmname)))
	return nil
}
