// vmman3 : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// vmmanagement/vmStop.go
// 2022-08-22 13:13:14

package vmmanagement

//import (
//	"fmt"
//	"libvirt.org/go/libvirt"
//	"os"
//	"vmman3/db"
//	"vmman3/domainmanagement"
//	"vmman3/helpers"
//	"vmman3/hypervisor"
//	"vmman3/inventory"
//)
//
//func Stop(args []string) error {
//	var bIsActive bool
//	var err error
//
//	conn := hypervisor.Connect2HVM()
//	defer conn.Close()
//
//	for _, vmname := range args {
//		var host string
//		domain := domainmanagement.GetDomain(conn, vmname)
//		if domain == nil {
//			os.Exit(0)
//		}
//		defer domain.Free()
//
//		bIsActive, _ = domain.IsActive()
//		if !bIsActive {
//			fmt.Printf("Domain %s on %s is already shut down\n", vmname, helpers.ConnectURI)
//		} else {
//			err = domain.ShutdownFlags(libvirt.DOMAIN_SHUTDOWN_DEFAULT)
//			fmt.Printf("Domain %s is being shut down... ", vmname)
//			if err != nil {
//				fmt.Printf("\nERROR: %s\n", helpers.Red(err.Error()))
//			} else {
//				fmt.Printf(helpers.Green("done\n"))
//				// This is where we update the vmstates table
//				if helpers.ConnectURI == "qemu:///system" {
//					host, _ = os.Hostname()
//				} else {
//					_, _, host = db.SplitConnectURI(helpers.ConnectURI)
//				}
//				vmStateChange(host, vmname)
//			}
//		}
//	}
//	return nil
//}
//
//func StopAll() error {
//	var vmlist []string
//	var domains []libvirt.Domain
//	var err error
//
//	if domains, err = inventory.GetVMlist(); err != nil {
//		return err
//	}
//	for _, domain := range domains {
//		_, err = domain.GetID()
//		if err == nil { // this means GetID() returned an ID, thus the VM is not shutdown (could be paused)
//			vmname, _ := domain.GetName()
//			vmlist = append(vmlist, vmname)
//		}
//	}
//	return Stop(vmlist)
//}
