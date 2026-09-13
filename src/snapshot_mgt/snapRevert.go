// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/snapshot_mgt/snapRevert.go
// Original timestamp : 2026.09.12 22:09:37

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

func RevertSnapshot(vmname, sname string) *ce.CustomError {
	var cerr *ce.CustomError
	var conn *libvirt.Connect
	var domain *libvirt.Domain

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

	if err = snapshot.RevertToSnapshot(0); err != nil {
		var lverr libvirt.Error
		if errors.As(err, &lverr) {
			return &ce.CustomError{Title: "Cannot revert " + vmname + " to snapshot " + sname, Message: lverr.Message}
		}
		return &ce.CustomError{Title: "Cannot revert " + vmname + " to snapshot " + sname, Message: err.Error()}
	}

	fmt.Println(hftx.EnabledSign(fmt.Sprintf("Snapshot %s %s from %s", sname,
		hftx.Bold(hftx.DimGreen("reverted")), vmname)))
	return nil
}
