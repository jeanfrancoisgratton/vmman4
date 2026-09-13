// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/snapshot_mgt/snapCreate.go
// Original timestamp : 2026.09.13 15:28:50

package snapshot_mgt

import (
	"encoding/xml"
	"errors"
	"fmt"
	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// domainSnapshotCreate : minimal XML payload accepted by Domain.CreateSnapshotXML()
type domainSnapshotCreate struct {
	XMLName     xml.Name `xml:"domainsnapshot"`
	Name        string   `xml:"name"`
	Description string   `xml:"description,omitempty"`
}

func CreateSnapshot(vmname, snapname, parentsnap, description string) *ce.CustomError {
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

	var origCurrent string
	if numsnap > 0 {
		if origCurrent, cerr = GetCurrentSnapshotName(conn, vmname); cerr != nil {
			return cerr
		}
		// no parent snapshot name was provided, so we'll branch off the current snapshot
		if parentsnap == "" {
			parentsnap = origCurrent
		}
	} else if parentsnap != "" {
		// no snapshot exists yet, so a specific parent can't possibly be found
		return &ce.CustomError{Title: "Snapshot not found",
			Message: "No snapshot named " + hftx.Bold(hftx.Red(parentsnap)) + " exists for " + hftx.Bold(hftx.Red(vmname))}
	}
	shared.Wait4Shutdown(domain, vmname)

	// branching off a snapshot other than the current one: check it out first
	if numsnap > 0 && parentsnap != origCurrent {
		if cerr = checkoutSnapshot(domain, vmname, parentsnap); cerr != nil {
			return cerr
		}
	}

	snapshotXML, err := xml.Marshal(domainSnapshotCreate{Name: snapname, Description: description})
	if err != nil {
		return &ce.CustomError{Title: "snapshot_mgt.CreateSnapshot() failure", Message: err.Error()}
	}

	newSnapshot, err := domain.CreateSnapshotXML(string(snapshotXML), 0)
	if err != nil {
		var lverr libvirt.Error
		if errors.As(err, &lverr) {
			return &ce.CustomError{Title: "Cannot create snapshot " + snapname + " on " + vmname, Message: lverr.Message}
		}
		return &ce.CustomError{Title: "Cannot create snapshot " + snapname + " on " + vmname, Message: err.Error()}
	}
	defer newSnapshot.Free()

	// restore the VM to whichever snapshot was current before this branch
	if numsnap > 0 && parentsnap != origCurrent {
		if cerr = checkoutSnapshot(domain, vmname, origCurrent); cerr != nil {
			return cerr
		}
	}

	fmt.Println(hftx.EnabledSign(fmt.Sprintf("Snapshot %s %s on %s", snapname,
		hftx.Bold(hftx.DimGreen("created")), vmname)))
	return nil
}
