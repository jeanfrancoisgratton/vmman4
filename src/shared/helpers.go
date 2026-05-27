// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/shared/helpers.go
// Original timestamp: 2026/05/27 08:46:14

package shared

import (
	"errors"
	"strings"

	"github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// Connect2HVM : this is where we actually handle the connection
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
	return conn, nil
}

// GetDomain : Connects to the VM, returning the domain object.
// Moved here from vmmanagement to break the inventory <-> vmmanagement import cycle.
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

// GetVMlist : Returns all domains (active + inactive) on the hypervisor.
// Moved here from inventory to break the vmmanagement -> inventory import cycle.
func GetVMlist() ([]libvirt.Domain, *customError.CustomError) {
	conn, cerr := Connect2HVM()
	if cerr != nil {
		return nil, cerr
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE | libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		return nil, &customError.CustomError{Title: "Error listing domains", Message: err.Error()}
	}
	return doms, nil
}
