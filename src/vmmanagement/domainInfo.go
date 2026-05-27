// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vmmanagement/domainInfo.go
// Original timestamp: 2026/05/26 19:05:31

package vmmanagement

import (
	"errors"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// GetDomain : Connects to the VM, returning the domain object
func GetDomain(conn *libvirt.Connect, vmname string) (*libvirt.Domain, *ce.CustomError) {
	domain, err := conn.LookupDomainByName(vmname)
	if err != nil {
		var lverr libvirt.Error
		ok := errors.As(err, &lverr)
		if ok {
			if strings.HasPrefix(lverr.Message, "Domain not found") {
				return nil, &ce.CustomError{Title: lverr.Message}
			} else {
				return nil, &ce.CustomError{Title: "domainmanagement.GetDomain() failure", Message: lverr.Message}
			}
		}
	}
	return domain, nil
}
