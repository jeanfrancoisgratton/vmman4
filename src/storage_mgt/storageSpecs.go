// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storage_mgt/storageSpecs.go
// Original timestamp: 2026/09/10

package storagemanagement

import (
	"encoding/xml"
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// GetStorageSpecs4VM returns the disk inventory (device, type, source, size) for the named domain.
// It parses the domain XML to enumerate disks and resolves each volume's logical
// size through the storage pool API where possible.
func GetStorageSpecs4VM(vmName string, conn *libvirt.Connect) (VMStorageInfo, *ce.CustomError) {
	var cerr *ce.CustomError
	//var conn *libvirt.Connect

	if conn == nil {
		if cerr = connection_mgt.ResolveConnectionURI(); cerr != nil {
			return VMStorageInfo{}, cerr
		}
		if conn, cerr = shared.Connect2HVM(); cerr != nil {
			return VMStorageInfo{}, cerr
		}
		defer conn.Close()
	}

	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return VMStorageInfo{}, &ce.CustomError{Title: "GetStorageSpecs4VM: domain lookup failed", Message: err.Error()}
	}
	defer dom.Free()

	xmlDesc, err := dom.GetXMLDesc(0)
	if err != nil {
		return VMStorageInfo{}, &ce.CustomError{Title: "GetStorageSpecs4VM: cannot get domain XML", Message: err.Error()}
	}

	var dXML domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &dXML); err != nil {
		return VMStorageInfo{}, &ce.CustomError{Title: "GetStorageSpecs4VM: cannot parse domain XML", Message: err.Error()}
	}

	info := VMStorageInfo{DomainName: vmName}

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

		di.SizeBytes, di.PoolName, di.PoolTargetPath, di.PoolState = resolveVolumeInfo(conn, di.SourcePath, d.Source.Pool, d.Source.Volume)
		info.Disks = append(info.Disks, di)
	}

	info.DiskCount = len(info.Disks)

	return info, nil
}
