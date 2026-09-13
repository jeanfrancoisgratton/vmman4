// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/snapshot_mgt/snapHelpers.go
// Original timestamp: 2026/05/26 19:33:08

package snapshot_mgt

import (
	"errors"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"

	"vmman4/shared"
)

// GetCurrentSnapshotName : Gets the name of the current snapshot for a given VM
func GetCurrentSnapshotName(conn *libvirt.Connect, vmname string) (string, *ce.CustomError) {
	var currentSnapshot = "n/a"
	var domain *libvirt.Domain
	var cerr *ce.CustomError
	var snapshots []libvirt.DomainSnapshot
	var err error

	if domain, cerr = shared.GetDomain(conn, vmname); cerr != nil {
		return "", cerr
	}

	defer domain.Free()

	if snapshots, err = domain.ListAllSnapshots(0); err != nil {
		return "", &ce.CustomError{Title: "Unable to list snapshots", Message: err.Error()}
	}

	for _, snapshot := range snapshots {
		var isCurrent, _ = snapshot.IsCurrent(0)
		if isCurrent {
			currentSnapshot, _ = snapshot.GetName()
			break
		}
	}
	return currentSnapshot, nil
}

// checkoutSnapshot : looks up sname on domain and reverts the domain to it
func checkoutSnapshot(domain *libvirt.Domain, vmname, sname string) *ce.CustomError {
	snapshot, err := domain.SnapshotLookupByName(sname, 0)
	if err != nil {
		var lverr libvirt.Error
		if errors.As(err, &lverr) {
			if strings.HasPrefix(lverr.Message, "Domain snapshot not found") {
				return &ce.CustomError{Title: "Snapshot not found",
					Message: "No snapshot named " + hftx.Bold(hftx.Red(sname)) + " exists for " + hftx.Bold(hftx.Red(vmname))}
			}
			return &ce.CustomError{Title: "snapshot_mgt.checkoutSnapshot() failure", Message: lverr.Message}
		}
		return &ce.CustomError{Title: "snapshot_mgt.checkoutSnapshot() failure", Message: err.Error()}
	}
	defer snapshot.Free()

	if err = snapshot.RevertToSnapshot(0); err != nil {
		var lverr libvirt.Error
		if errors.As(err, &lverr) {
			return &ce.CustomError{Title: "Cannot revert " + vmname + " to snapshot " + sname, Message: lverr.Message}
		}
		return &ce.CustomError{Title: "Cannot revert " + vmname + " to snapshot " + sname, Message: err.Error()}
	}
	return nil
}
