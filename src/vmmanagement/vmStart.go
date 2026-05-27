// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vmmanagement/vmStart.go
// Original timestamp: 2026/05/27 08:46:14

package vmmanagement

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
	"vmman4/connection"
	"vmman4/shared"
)

// Start : starts one or many VMs
func Start(args []string) *ce.CustomError {
	var bIsActive bool
	var err *ce.CustomError
	var conn *libvirt.Connect
	var domain *libvirt.Domain

	if err = connection.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = shared.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	for _, vmname := range args {
		if domain, err = shared.GetDomain(conn, vmname); err != nil {
			return err
		}
		defer domain.Free()

		bIsActive, _ = domain.IsActive()
		if bIsActive {
			fmt.Println(hftx.WarningSign("Domain " + vmname + " already active"))
		} else {
			fmt.Printf("%s", hftx.InProgressSign("Domain "+vmname+" is starting... "))
			if serr := domain.Create(); serr != nil {
				fmt.Println()
				return &ce.CustomError{Title: "Could not start " + vmname, Message: serr.Error()}
			} else {
				fmt.Println(hftx.Green("DONE"))
			}
		}
	}
	return nil
}

// StartAll : fetches the list of VMs on the hypervisor, and then starts them
func StartAll() *ce.CustomError {
	var vmlist []string
	var domains []libvirt.Domain
	var err *ce.CustomError

	if err = connection.ResolveConnectionURI(); err != nil {
		return err
	}
	if domains, err = shared.GetVMlist(); err != nil {
		return err
	}

	for _, domain := range domains {
		_, serr := domain.GetID()
		if serr != nil { // GetID() failed → domain has no ID → it is not running; candidate to start
			vmname, _ := domain.GetName()
			vmlist = append(vmlist, vmname)
		}
	}
	return Start(vmlist)
}
