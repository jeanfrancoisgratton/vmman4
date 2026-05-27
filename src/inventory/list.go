// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/inventory/list.go
// Original timestamp: 2026/05/22 07:54:46

package inventory

import (
	"fmt"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"libvirt.org/go/libvirt"
	"vmman4/connection"
	"vmman4/shared"
)

func VmInventory() *ce.CustomError {
	var (
		vmspecs []vmInfo
		conn    *libvirt.Connect
		err     *ce.CustomError
	)

	// 1st step : connect to the hypervisor, wherever it is (-c or -C flag)
	if err = connection.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = connection.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	// 2nd step: collect the information
	if vmspecs, err = collectInfo(conn); err != nil {
		return err
	}

	// 3rd step: display information
	fmt.Println("Domains on hypervisor " + hftx.Blue(shared.ConnectURI) + ":")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "VM name", "State", "vMem", "vCPUs", "Snaps", "Curr snap", "IP", "Last seen / uptime", "Hypervisor", "OS", "Storage"})

	for _, vmspec := range vmspecs {
		sID := ""
		if vmspec.viId > 0 && vmspec.viId < 10 {
			sID = fmt.Sprintf("000%d", vmspec.viId)
		}
		if vmspec.viId > 9 && vmspec.viId < 100 {
			sID = fmt.Sprintf("00%d", vmspec.viId)
		}
		if vmspec.viId > 99 && vmspec.viId < 999 {
			sID = fmt.Sprintf("0%d", vmspec.viId)
		}
		// time.Unix(image.Created, 0).Format("2006.01.02 15:04:05")
		t.AppendRow([]interface{}{sID, vmspec.viName, vmspec.viState, vmspec.viMem, vmspec.viCpu})
	}
	t.SortBy([]table.SortBy{
		{Name: "ID", Mode: table.Asc},
		{Name: "VM name", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	//t.Style().Options.DrawBorder = false
	//t.Style().Options.SeparateColumns = false
	t.Style().Format.Header = text.FormatDefault
	t.SetRowPainter(func(row table.Row) text.Colors {
		switch row[2] {
		case "Running":
			//return text.Colors{text.BgBlack, text.FgHiGreen}
			return text.Colors{text.FgHiGreen}
		case "Crashed":
			//return text.Colors{text.BgBlack, text.FgHiRed}
			return text.Colors{text.FgHiRed}
		case "Blocked":
		case "Suspended":
		case "Paused":
			//return text.Colors{text.BgHiBlack, text.FgHiYellow}
			return text.Colors{text.FgHiYellow}
		}
		return nil
	})
	t.Render()
	return nil
}
