// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storageManagement/storageList.go
// Original timestamp: 2026/05/30 14:09:11

package storagemanagement

import (
	"encoding/xml"
	"fmt"
	"os"

	"vmman4/connection"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"libvirt.org/go/libvirt"
)

// ListVMStorage returns the number of disks and their details for the named domain.
// It parses the domain XML to enumerate disks and resolves each volume's logical
// size through the storage pool API where possible.
func ListVMStorage(domainName string, displayOutput bool) (VMStorageInfo, *ce.CustomError) {
	var cerr *ce.CustomError
	var conn *libvirt.Connect

	if cerr = connection.ResolveConnectionURI(); cerr != nil {
		return VMStorageInfo{}, cerr
	}
	if conn, cerr = shared.Connect2HVM(); cerr != nil {
		return VMStorageInfo{}, cerr
	}
	defer conn.Close()

	dom, err := conn.LookupDomainByName(domainName)
	if err != nil {
		return VMStorageInfo{}, &ce.CustomError{Title: "ListVMStorage: domain lookup failed", Message: err.Error()}
	}
	defer dom.Free()

	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		return VMStorageInfo{}, &ce.CustomError{Title: "ListVMStorage: cannot get domain XML", Message: err.Error()}
	}

	var dXML domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &dXML); err != nil {
		return VMStorageInfo{}, &ce.CustomError{Title: "ListVMStorage: cannot parse domain XML", Message: err.Error()}
	}

	info := VMStorageInfo{DomainName: domainName}

	for _, d := range dXML.Devices.Disks {
		if d.Device == "cdrom" || d.Device == "floppy" {
			continue
		}

		di := DiskInfo{Device: d.Target.Dev, Type: d.Type}

		switch d.Type {
		case "file":
			di.SourcePath = d.Source.File
		case "block":
			di.SourcePath = d.Source.Dev
		case "network":
			di.SourcePath = d.Source.Name
		default:
			if d.Source.Pool != "" && d.Source.Volume != "" {
				di.SourcePath = fmt.Sprintf("pool:%s/volume:%s", d.Source.Pool, d.Source.Volume)
			}
		}

		di.SizeBytes = resolveVolumeSize(conn, di.SourcePath, d.Source.Pool, d.Source.Volume)
		info.Disks = append(info.Disks, di)
	}

	info.DiskCount = len(info.Disks)

	if !displayOutput {
		return info, nil
	}

	// displayOutput is set to true, we render the received output
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Disk", "Device", "Disk type", "Source path", "Size (GB)"})
	var vSizeString string
	var berr error

	for ndx, di := range info.Disks {
		if vSizeString, berr = hf.BytesToUnit(di.SizeBytes, 'g', 3); berr != nil {
			return VMStorageInfo{}, &ce.CustomError{Title: "ListVMStorage: size conversion failed",
				Message: berr.Error()}
		}
		t.AppendRow([]interface{}{ndx, "/dev/" + di.Device, di.Type, di.SourcePath, vSizeString + " GB"})

	}

	t.SortBy([]table.SortBy{
		{Name: "Disk", Mode: table.Asc},
	})
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	fmt.Println()

	return VMStorageInfo{}, nil
}
