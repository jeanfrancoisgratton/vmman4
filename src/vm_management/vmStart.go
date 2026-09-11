// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_management/vmStart.go
// Original timestamp: 2026/05/27 08:46:14

package vm_management

import (
	"fmt"

	"vmman4/connection"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// Start : starts one or many VMs
func Start(args []string) *ce.CustomError {
	if err := connection.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	return startDomains(conn, args)
}

// startDomains : starts one or many VMs using an already-open connection
func startDomains(conn *libvirt.Connect, args []string) *ce.CustomError {
	var bIsActive bool
	var err *ce.CustomError
	var domain *libvirt.Domain

	for _, vmname := range args {
		if domain, err = shared.GetDomain(conn, vmname); err != nil {
			return err
		}
		defer domain.Free()

		bIsActive, _ = domain.IsActive()
		if bIsActive {
			fmt.Println(hftx.WarningSign("Domain " + vmname + " already active"))
		} else {
			if serr := domain.Create(); serr != nil {
				fmt.Println()
				return &ce.CustomError{Title: "Could not start " + vmname, Message: serr.Error()}
			} else {
				fmt.Println(hftx.EnabledSign("Domain " + vmname + hftx.Green(" started")))
			}
		}
	}
	return nil
}

// StartAll : fetches the list of VMs on the hypervisor, and then starts them
func StartAll() *ce.CustomError {
	if err := connection.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	domains, serr := listDomains(conn)
	if serr != nil {
		return serr
	}

	var vmlist []string
	for _, domain := range domains {
		_, e := domain.GetID()
		if e != nil { // GetID() failed → domain has no ID → it is not running; candidate to start
			vmname, _ := domain.GetName()
			vmlist = append(vmlist, vmname)
		}
	}
	return startDomains(conn, vmlist)
}
