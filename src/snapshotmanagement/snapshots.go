// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/snapshotmanagement/snapshots.go
// Original timestamp: 2026/05/26 19:33:08

package snapshotmanagement

import (
	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
	"vmman4/vmmanagement"
)

// GetCurrentSnapshotName : Gets the name of the current snapshot for a given VM
func GetCurrentSnapshotName(conn *libvirt.Connect, vmname string) (string, *ce.CustomError) {
	var currentSnapshot = "n/a"
	var domain *libvirt.Domain
	var cerr *ce.CustomError
	var snapshots []libvirt.DomainSnapshot
	var err error

	if domain, cerr = vmmanagement.GetDomain(conn, vmname); cerr != nil {
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
