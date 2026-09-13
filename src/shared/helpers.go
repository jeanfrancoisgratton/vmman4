// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/shared/helpers.go
// Original timestamp: 2026/05/27 08:46:14

package shared

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// Connect2HVM : this is where we actually handle the connection_mgt
func Connect2HVM() (*libvirt.Connect, *customError.CustomError) {
	conn, err := libvirt.NewConnect(ConnectURI)
	if err != nil {
		var lverr libvirt.Error
		ok := errors.As(err, &lverr)
		if ok && (lverr.Message == "End of file while reading data: virt-ssh-helper: cannot connect to '/var/run/libvirt/libvirt-sock': Failed to connect socket to '/var/run/libvirt/libvirt-sock': Connection refused: Input/output error" ||
			strings.HasPrefix(lverr.Message, "Cannot recv data: ssh: connect to host")) {
			return nil, &customError.CustomError{Title: "Unable to connect to host", Message: "Hypervisor " + ConnectURI + " is offline"}
		} else {
			return nil, &customError.CustomError{Title: "Unable to connect to host", Message: err.Error()}
		}
	}
	if !QuietOutput {
		fmt.Println()
		fmt.Println(hftx.InfoSign(" Connected on hypervisor " + hftx.Blue(ConnectURI)))
		fmt.Println()
	}
	return conn, nil
}

// GetDomain : Connects to the VM, returning the domain object.
// Moved here from vm_mgt to break the inventory <-> vm_mgt import cycle.
func GetDomain(conn *libvirt.Connect, vmname string) (*libvirt.Domain, *customError.CustomError) {
	domain, err := conn.LookupDomainByName(vmname)
	if err != nil {
		var lverr libvirt.Error
		ok := errors.As(err, &lverr)
		if ok {
			if strings.HasPrefix(lverr.Message, "Domain not found") {
				return nil, &customError.CustomError{Title: lverr.Message}
			}
			return nil, &customError.CustomError{Title: "shared.GetDomain() failure", Message: lverr.Message}
		}
	}
	return domain, nil
}

// Wait4Shutdown : Tries 15 seconds to gracefully shutdown the VM, if not it will shutdown forcefully
func Wait4Shutdown(vm *libvirt.Domain, vmname string) {
	var bIsActive = false
	fmt.Println(hftx.InProgressSign("Waiting for VM to shutdown..."))
	//fmt.Println("Will await that the VM " + vmname + " gracefully shuts down on " + shared.ConnectURI)
	bIsActive, _ = vm.IsActive()
	if bIsActive {
		n := 15
		vm.DestroyFlags(libvirt.DOMAIN_DESTROY_GRACEFUL)
		for n > 0 {
			bIsActive, _ = vm.IsActive()
			if bIsActive {
				n -= 1
				time.Sleep(1 * time.Second)
			} else {
				n = 0
			}
		}
		bIsActive, _ = vm.IsActive()
		if bIsActive {
			vm.DestroyFlags(libvirt.DOMAIN_DESTROY_DEFAULT)
			fmt.Println("The VM " + vmname + " was slow to shutdown and was forcely shut down")
		}
	}
}
