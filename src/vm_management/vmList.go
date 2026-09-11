// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/inventory/vmList.go
// Original timestamp: 2026/05/22 07:54:46

package vm_management

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"vmman4/connection"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"libvirt.org/go/libvirt"
)

func VmInventory() *ce.CustomError {
	var (
		vmspecs []vmInfo
		conn    *libvirt.Connect
		err     *ce.CustomError
	)

	if err = connection.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = shared.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	if vmspecs, err = collectInfo(conn); err != nil {
		return err
	}

	hypervisor := ""
	if strings.HasPrefix(shared.ConnectURI, "qemu:///") {
		hypervisor = "local"
	} else {
		n := strings.Index(shared.ConnectURI, "@") + 1
		hypervisor = strings.TrimSuffix(shared.ConnectURI[n:], "/system")
	}

	//fmt.Println()
	fmt.Println(fmt.Sprintf("\n%s domains found on hypervisor %s\n",
		hftx.Green(strconv.Itoa(len(vmspecs))), hftx.Green(hypervisor)))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "VM name", "State", "vMem", "vCPUs", "Snapshots", "Current snapshot", "Interface", "IP addr"})

	for _, vmspec := range vmspecs {
		sID := fmt.Sprintf("%04d", vmspec.viId)
		t.AppendRow([]interface{}{sID, vmspec.viName, vmspec.viState, vmspec.viMem, vmspec.viCpu,
			vmspec.viSnapshots, vmspec.viCurrentSnapshot, vmspec.viInterfaceName, vmspec.viIPaddress})
	}

	t.SortBy([]table.SortBy{
		{Name: "ID", Mode: table.Asc},
		{Name: "VM name", Mode: table.Asc},
	})

	//vmlistStyle := table.StyleColoredBlackOnYellowWhite
	vmlistStyle := table.StyleRounded
	//vmlistStyle.Color.Row = text.Colors{text.FgHiWhite}
	//vmlistStyle.Color.RowAlternate = text.Colors{text.FgHiWhite, text.BgBlack}
	t.SetStyle(vmlistStyle)
	t.Style().Format.Header = text.FormatDefault

	// Stretch table to terminal width: cap at terminal width; "VM name"
	// absorbs any leftover space since it's the most variable column.
	tw := termWidth()
	t.SetAllowedRowLength(tw)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "VM name", WidthMin: 16, WidthMax: tw / 3},
	})

	// SetRowPainter overrides Row/RowAlternate only when it returns non-nil.
	// Returning nil falls back to the style's alternating colours.
	t.SetRowPainter(func(row table.Row) text.Colors {
		switch row[2] {
		case "Running":
			return text.Colors{text.FgHiGreen}
		case "Crashed":
			return text.Colors{text.FgHiRed}
		case "Blocked", "Suspended", "Paused":
			return text.Colors{text.FgHiYellow}
		}
		return nil
	})

	t.Render()
	return nil
}

// GetVMlist : Returns all domains (active + inactive) on the hypervisor.
// Moved here from inventory to break the vm_management -> inventory import cycle.
func GetVMlist() ([]libvirt.Domain, *ce.CustomError) {
	if cerr := connection.ResolveConnectionURI(); cerr != nil {
		return nil, cerr
	}
	conn, cerr := shared.Connect2HVM()
	if cerr != nil {
		return nil, cerr
	}
	defer conn.Close()

	return listDomains(conn)
}

// listDomains : Returns all domains (active + inactive) using an already-open connection.
func listDomains(conn *libvirt.Connect) ([]libvirt.Domain, *ce.CustomError) {
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return nil, &ce.CustomError{Title: "Error listing domains", Message: err.Error()}
	}
	return doms, nil
}
