// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/snapshot_mgt/snapDumpXML.go
// Original timestamp : 2026.09.13 16:00:00

package snapshot_mgt

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// DumpSnapshotXML : dumps a snapshot's XML description, to stdout or to filename if provided
func DumpSnapshotXML(vmname, sname, filename string) *ce.CustomError {
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
	if numsnap == 0 {
		return &ce.CustomError{Fatality: ce.NotAnError, Message: "Snapshot number is zero"}
	}
	// no snapshot name were provided, we'll use the current snapshot
	if sname == "" {
		if sname, cerr = GetCurrentSnapshotName(conn, vmname); cerr != nil {
			return cerr
		}
	}

	snapshot, err := domain.SnapshotLookupByName(sname, 0)
	if err != nil {
		var lverr libvirt.Error
		if errors.As(err, &lverr) {
			if strings.HasPrefix(lverr.Message, "Domain snapshot not found") {
				return &ce.CustomError{Title: "Snapshot not found",
					Message: "No snapshot named " + hftx.Bold(hftx.Red(sname)) + " exists for " + hftx.Bold(hftx.Red(vmname))}
			}
			return &ce.CustomError{Title: "snapshot_mgt.DumpSnapshotXML() failure", Message: lverr.Message}
		}
		return &ce.CustomError{Title: "snapshot_mgt.DumpSnapshotXML() failure", Message: err.Error()}
	}
	defer snapshot.Free()

	data, err := snapshot.GetXMLDesc(0)
	if err != nil {
		return &ce.CustomError{Title: "Cannot dump XML for snapshot " + sname, Message: err.Error()}
	}

	if filename == "" {
		fmt.Println(data)
		return nil
	}

	if !strings.HasSuffix(filename, ".xml") {
		filename += ".xml"
	}
	file, e := os.Create(filename)
	if e != nil {
		return &ce.CustomError{Title: "os.Create() failed", Message: e.Error()}
	}
	defer file.Close()
	if _, e = file.WriteString(data); e != nil {
		return &ce.CustomError{Title: "file.WriteString() : Unable to write XML file", Message: e.Error()}
	}
	file.Sync()

	fmt.Println(hftx.EnabledSign(fmt.Sprintf("Snapshot %s XML dumped to %s", sname, filename)))
	return nil
}
