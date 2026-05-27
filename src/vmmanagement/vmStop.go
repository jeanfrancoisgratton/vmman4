// vmman3 : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// vmmanagement/vmStop.go
// 2022-08-22 13:13:14

package vmmanagement

import (
	"fmt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
	"vmman4/connection"
	"vmman4/shared"
)

// Stop : stops one or many VMs
func Stop(args []string) *ce.CustomError {
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

		var e error
		bIsActive, e = domain.IsActive()
		if e != nil {
			fmt.Println(e.Error())
		}
		if !bIsActive {
			fmt.Println(hftx.WarningSign("Domain " + vmname + " is already shut down"))
		} else {
			fmt.Printf("%s", hftx.InProgressSign("Domain "+vmname+" is shutting down... "))
			if serr := domain.ShutdownFlags(libvirt.DOMAIN_SHUTDOWN_DEFAULT); serr != nil {
				fmt.Println()
				return &ce.CustomError{Title: "Could not stop " + vmname, Message: serr.Error()}
			} else {
				fmt.Println(hftx.Green("DONE"))
			}
		}
	}
	return nil
}

// StopAll : fetches the list of VMs on the hypervisor and then stops them
func StopAll() *ce.CustomError {
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
		if serr == nil { // GetID() failed → domain has no ID → it is not running; candidate to start
			vmname, _ := domain.GetName()
			vmlist = append(vmlist, vmname)
		}
	}
	return Stop(vmlist)
}
