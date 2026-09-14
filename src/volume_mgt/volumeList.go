// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/volumeList.go

package volume_mgt

import (
	"fmt"
	"os"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"libvirt.org/go/libvirt"
)

// volumeRow is a single line of the `vol list` table.
type volumeRow struct {
	pool string
	name string
	path string
	size uint64
}

// ListVolumes prints all volumes in the named storage pool, or across every
// active pool on the hypervisor when poolName is empty.
func ListVolumes(poolName string) *ce.CustomError {
	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	var pools []libvirt.StoragePool
	if poolName != "" {
		pool, perr := conn.LookupStoragePoolByName(poolName)
		if perr != nil {
			return &ce.CustomError{Title: fmt.Sprintf("ListVolumes: pool %q not found", poolName), Message: perr.Error()}
		}
		pools = append(pools, *pool)
	} else {
		active, lerr := conn.ListAllStoragePools(libvirt.CONNECT_LIST_STORAGE_POOLS_ACTIVE)
		if lerr != nil {
			return &ce.CustomError{Title: "ListVolumes: cannot list active pools", Message: lerr.Error()}
		}
		pools = active
	}

	var rows []volumeRow
	for _, p := range pools {
		defer p.Free()
		pname, _ := p.GetName()

		vols, verr := p.ListAllStorageVolumes(0)
		if verr != nil {
			continue
		}
		for _, v := range vols {
			vname, _ := v.GetName()
			vpath, _ := v.GetPath()
			var vsize uint64
			if vi, ierr := v.GetInfo(); ierr == nil {
				vsize = vi.Capacity
			}
			rows = append(rows, volumeRow{pool: pname, name: vname, path: vpath, size: vsize})
			v.Free()
		}
	}

	if shared.QuietOutput {
		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Volume name", "Pool", "Size (GB)", "Path"})

	for _, r := range rows {
		sizeGB, cerr := hf.BytesToUnit(r.size, 'g', 3)
		if cerr != nil {
			return &ce.CustomError{Title: "Unable to convert disk size units", Message: cerr.Error()}
		}
		t.AppendRow([]interface{}{r.name, r.pool, hf.SI(sizeGB) + " GB", r.path})
	}

	t.SortBy([]table.SortBy{
		{Name: "Pool", Mode: table.Asc},
		{Name: "Volume name", Mode: table.Asc},
	})
	t.SetStyle(table.StyleRounded)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}
