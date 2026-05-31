// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storageManagement/storageList.go
// Original timestamp: 2026/05/30 14:09:11

package storagemanagement

import (
	"encoding/xml"
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// ListVMStorage returns the number of disks and their details for the named domain.
// It parses the domain XML to enumerate disks and resolves each volume's logical
// size through the storage pool API where possible.
func ListVMStorage(conn *libvirt.Connect, domainName string) (*VMStorageInfo, *ce.CustomError) {
	dom, err := conn.LookupDomainByName(domainName)
	if err != nil {
		return nil, &ce.CustomError{Title: "ListVMStorage: domain lookup failed", Message: err.Error()}
	}
	defer dom.Free()

	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		return nil, &ce.CustomError{Title: "ListVMStorage: cannot get domain XML", Message: err.Error()}
	}

	var dXML domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &dXML); err != nil {
		return nil, &ce.CustomError{Title: "ListVMStorage: cannot parse domain XML", Message: err.Error()}
	}

	info := &VMStorageInfo{DomainName: domainName}

	for _, d := range dXML.Devices.Disks {
		if d.Device == "cdrom" || d.Device == "floppy" {
			continue
		}

		di := DiskInfo{
			Device: d.Target.Dev,
			Type:   d.Type,
		}

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
	return info, nil
}
