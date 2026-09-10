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
	if err := connection.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	return stopDomains(conn, args)
}

// stopDomains : stops one or many VMs using an already-open connection
func stopDomains(conn *libvirt.Connect, args []string) *ce.CustomError {
	var bIsActive bool
	var err *ce.CustomError
	var domain *libvirt.Domain

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
			if serr := domain.ShutdownFlags(libvirt.DOMAIN_SHUTDOWN_DEFAULT); serr != nil {
				fmt.Println()
				return &ce.CustomError{Title: "Could not stop " + vmname, Message: serr.Error()}
			} else {
				fmt.Println(hftx.EnabledSign("Domain " + vmname + hftx.Red(" stopped")))
			}
		}
	}
	return nil
}

// StopAll : fetches the list of VMs on the hypervisor and then stops them
func StopAll() *ce.CustomError {
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
		if e == nil { // GetID() succeeded → domain has an ID → it is running; candidate to stop
			vmname, _ := domain.GetName()
			vmlist = append(vmlist, vmname)
		}
	}
	return stopDomains(conn, vmlist)
}
