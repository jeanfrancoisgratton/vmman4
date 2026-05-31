// vmman4
// Written by J.F. Gratton (jean-francois@famillegratton.net)
// Original filename : src/snapshotmanagement/snapList.go
// original timestamp : 2026/05/31 14:57:46

package snapshotmanagement

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
	"vmman4/connection"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"libvirt.org/go/libvirt"
)

// ListSnapshots() : Lists all snapshots on current VM
func ListSnapshots(vmname string) *ce.CustomError {
	var snapXMLdata SnapshotXMLstruct
	var snaps []SnapshotXMLstruct
	var cerr *ce.CustomError
	var conn *libvirt.Connect
	var domain *libvirt.Domain

	if cerr = connection.ResolveConnectionURI(); cerr != nil {
		return cerr
	}
	if conn, cerr = shared.Connect2HVM(); cerr != nil {
		return cerr
	}
	defer conn.Close()

	if domain, cerr = shared.GetDomain(conn, vmname); cerr != nil {
		return cerr
	}
	if domain == nil {
		os.Exit(0)
	}
	defer domain.Free()
	numsnap, _ := domain.SnapshotNum(0)
	if numsnap == 0 {
		fmt.Printf("Domain %s has no snapshot\n", vmname)
		os.Exit(0)
	}
	snapshots, _ := domain.ListAllSnapshots(0)

	for _, snap := range snapshots {
		data, _ := snap.GetXMLDesc(0)
		if err := xml.Unmarshal([]byte(data), &snapXMLdata); err != nil {
			fmt.Println("err:", err)
			os.Exit(-1)
		} else {

			snapXMLdata.CurrentSnapshot, _ = snap.IsCurrent(0)
			snaps = append(snaps, snapXMLdata)
		}
	}
	displaySnapshots(snaps, vmname)
	return nil
}

// displaySnapshots() : will display the actual snapshot info in a table
func displaySnapshots(snaps []SnapshotXMLstruct, vmname string) {

	//helpers.SurroundText(fmt.Sprintf("All snapshots on %s/%s", helpers.ConnectURI, vmname), false)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Snapshot name", "Current", "Parent", "Creation time"})

	for _, snapshot := range snaps {
		tcreated := time.Unix(snapshot.CreationTime, 0).Format("2006.01.02 15:04:05")
		t.AppendRow([]interface{}{snapshot.SnapshotName, snapshot.CurrentSnapshot, snapshot.Parent.ParentName, tcreated})

	}
	t.SortBy([]table.SortBy{
		{Name: "Snapshot name", Mode: table.Asc},
	})
	t.SetStyle(table.StyleDefault)
	t.Style().Options.DrawBorder = false
	//t.Style().Options.SeparateColumns = false
	t.Style().Format.Header = text.FormatDefault
	t.SetRowPainter(func(row table.Row) text.Colors {
		switch row[1] {
		case true:
			return text.Colors{text.FgHiGreen}
		}
		return nil
	})
	t.Render()
}
