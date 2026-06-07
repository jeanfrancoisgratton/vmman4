// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storageManagement/poolList.go
// Original timestamp: 2026/05/30 14:09:25

package storagemanagement

import (
	"encoding/xml"
	"os"

	"vmman4/connection"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"libvirt.org/go/libvirt"
)

// ListStoragePools enumerates all storage pools (active and inactive) and
// returns their paths / disk names along with the volumes each contains.
func ListStoragePools() ([]StoragePoolInfo, *ce.CustomError) {
	var cerr *ce.CustomError
	var conn *libvirt.Connect

	if cerr = connection.ResolveConnectionURI(); cerr != nil {
		return nil, cerr
	}
	if conn, cerr = shared.Connect2HVM(); cerr != nil {
		return nil, cerr
	}
	defer conn.Close()

	activePools, err := conn.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE)
	if err != nil {
		return nil, &ce.CustomError{Title: "ListStoragePools: cannot list active pools", Message: err.Error()}
	}

	inactivePools, err := conn.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_INACTIVE)
	if err != nil {
		return nil, &ce.CustomError{Title: "ListStoragePools: cannot list inactive pools", Message: err.Error()}
	}

	allPools := append(activePools, inactivePools...)
	result := make([]StoragePoolInfo, 0, len(allPools))

	for _, p := range allPools {
		defer p.Free()

		name, err := p.GetName()
		if err != nil {
			continue
		}

		uuid, _ := p.GetUUIDString()
		pinfo := StoragePoolInfo{Name: name, UUID: uuid}

		if poolInfo, err := p.GetInfo(); err == nil {
			switch poolInfo.State {
			case libvirt.STORAGE_POOL_RUNNING:
				pinfo.State = "running"
			case libvirt.STORAGE_POOL_INACTIVE:
				pinfo.State = "inactive"
			case libvirt.STORAGE_POOL_BUILDING:
				pinfo.State = "building"
			case libvirt.STORAGE_POOL_DEGRADED:
				pinfo.State = "degraded"
			case libvirt.STORAGE_POOL_INACCESSIBLE:
				pinfo.State = "inaccessible"
			default:
				pinfo.State = "unknown"
			}
		}

		if xmlDesc, err := p.GetXMLDesc(0); err == nil {
			var px poolXML
			if xml.Unmarshal([]byte(xmlDesc), &px) == nil {
				pinfo.TargetPath = px.Target.Path
			}
		}

		if pinfo.State == "running" {
			if vols, err := p.ListAllStorageVolumes(0); err == nil {
				for _, v := range vols {
					defer v.Free()
					vname, _ := v.GetName()
					vpath, _ := v.GetPath()
					var vcap uint64
					if vi, err := v.GetInfo(); err == nil {
						vcap = vi.Capacity
					}
					pinfo.Volumes = append(pinfo.Volumes, VolumeInfo{
						Name:      vname,
						Path:      vpath,
						SizeBytes: vcap,
					})
				}
			}
		}

		result = append(result, pinfo)
	}

	if !shared.QuietOutput {
		if err := displayPoolList(result); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func displayPoolList(result []StoragePoolInfo) *ce.CustomError {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Pool name", "UUID", "# volumes",
		"Pool path", "Pool state", "Pool size (GB)"})
	for _, p := range result {
		var vSizes string
		var volumeSize uint64
		for _, v := range p.Volumes {
			volumeSize += v.SizeBytes
		}
		if convertedSize, e := hf.BytesToUnit(volumeSize, 'g', 3); e != nil {
			return &ce.CustomError{Title: "Unable to convert disk size units", Message: e.Error()}
		} else {
			vSizes = hf.SI(convertedSize) + " GB"
		}

		t.AppendRow([]interface{}{p.Name, p.UUID, len(p.Volumes), p.TargetPath, p.State, vSizes})
	}
	t.SortBy([]table.SortBy{
		{Name: "Pool name", Mode: table.Asc},
	})

	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}
