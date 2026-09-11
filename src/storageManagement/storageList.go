// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/storageManagement/storageList.go
// Original timestamp: 2026/05/30 14:09:11

package storagemanagement

import (
	"fmt"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ListVMStorage returns the number of disks and their details for the named domain.
// It delegates disk enumeration to GetStorageSpecs4VM and optionally renders the result.
func ListVMStorage(domainName string, displayOutput bool) (VMStorageInfo, *ce.CustomError) {
	info, cerr := GetStorageSpecs4VM(domainName, nil)
	if cerr != nil {
		return VMStorageInfo{}, cerr
	}

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
