// vmman4
// Written by J.F. Gratton (jean-francois@famillegratton.net)
// Original filename : src/snapshot_mgt/snapList.go
// original timestamp : 2026/05/31 14:57:46

package snapshot_mgt

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"time"
	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"libvirt.org/go/libvirt"
)

// ListSnapshots() : Lists all snapshots on named VM(s). If asTree is set, the snapshot
// hierarchy is rendered as a tree instead of a flat table.

func ListSnapshots(vmnames []string, asTree bool) *ce.CustomError {
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

	fmt.Println()
	for _, vmname := range vmnames {
		if domain, cerr = shared.GetDomain(conn, vmname); cerr != nil {
			return cerr
		}
		defer domain.Free()
		numsnap, _ := domain.SnapshotNum(0)
		if numsnap == 0 {
			hftx.InfoSign("No snapshot on domain " + vmname)
		}
		snapshots, _ := domain.ListAllSnapshots(0)

		var snaps []SnapshotXMLstruct
		for _, snap := range snapshots {
			var snapXMLdata SnapshotXMLstruct
			data, _ := snap.GetXMLDesc(0)
			if err := xml.Unmarshal([]byte(data), &snapXMLdata); err != nil {
				return &ce.CustomError{Title: "Error unmarshalling data", Message: err.Error()}
			} else {
				snapXMLdata.CurrentSnapshot, _ = snap.IsCurrent(0)
				snaps = append(snaps, snapXMLdata)
			}
		}
		if asTree {
			displaySnapshotsTree(snaps, vmname)
		} else {
			displaySnapshots(snaps, vmname)
		}
	}
	return nil
}

// displaySnapshots() : will display the actual snapshot info in a table
func displaySnapshots(snaps []SnapshotXMLstruct, vmname string) {

	fmt.Println("Snapshot information for: " + hftx.Bold(hftx.White(vmname)))
	fmt.Println()
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

	fmt.Println()
}

// snapTreeNode : one node of the parent/child snapshot hierarchy
type snapTreeNode struct {
	name     string
	current  bool
	children []*snapTreeNode
}

// buildSnapshotForest : arranges snaps into a forest of snapTreeNode based on their parent name
func buildSnapshotForest(snaps []SnapshotXMLstruct) []*snapTreeNode {
	nodes := make(map[string]*snapTreeNode, len(snaps))
	for _, s := range snaps {
		nodes[s.SnapshotName] = &snapTreeNode{name: s.SnapshotName, current: s.CurrentSnapshot}
	}

	var roots []*snapTreeNode
	for _, s := range snaps {
		node := nodes[s.SnapshotName]
		if parent, ok := nodes[s.Parent.ParentName]; ok {
			parent.children = append(parent.children, node)
		} else {
			roots = append(roots, node)
		}
	}

	byName := func(n []*snapTreeNode) { sort.Slice(n, func(i, j int) bool { return n[i].name < n[j].name }) }
	byName(roots)
	for _, n := range nodes {
		byName(n.children)
	}
	return roots
}

// renderSnapshotTree : prints nodes as a `tree`-style box-drawing hierarchy
func renderSnapshotTree(nodes []*snapTreeNode, prefix string) {
	for i, n := range nodes {
		last := i == len(nodes)-1
		connector, childPrefix := "├── ", prefix+"│   "
		if last {
			connector, childPrefix = "└── ", prefix+"    "
		}
		label := n.name
		if n.current {
			label = hftx.Bold(hftx.DimGreen(n.name)) + " (current)"
		}
		fmt.Println(prefix + connector + label)
		renderSnapshotTree(n.children, childPrefix)
	}
}

// displaySnapshotsTree() : will display the snapshot hierarchy for vmname as a tree
func displaySnapshotsTree(snaps []SnapshotXMLstruct, vmname string) {
	fmt.Println("Snapshot tree for: " + hftx.Bold(hftx.White(vmname)))
	fmt.Println()
	if len(snaps) == 0 {
		fmt.Println()
		return
	}
	renderSnapshotTree(buildSnapshotForest(snaps), "")
	fmt.Println()
}
