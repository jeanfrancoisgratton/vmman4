// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vmmanagement/vmHelpers.go
// Original timestamp: 2026/05/26 19:08:41

package vmmanagement

import (
	"fmt"
	"time"

	"libvirt.org/go/libvirt"
	"vmman4/shared"
)

// Wait4Shutdown ; Tries 15 seconds to gracefully shutdown the VM, if not it will shutdown forcefully
func Wait4Shutdown(vm *libvirt.Domain, vmname string) {
	var bIsActive = false
	fmt.Println("Will await that the VM " + vmname + " gracefully shuts down on " + shared.ConnectURI)
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
