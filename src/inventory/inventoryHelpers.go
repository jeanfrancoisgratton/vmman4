// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/inventory/inventoryHelpers.go
// Original timestamp: 2026/05/26 17:25:00

package inventory

import (
	"strconv"
	"time"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// getStateHelper : gets DomainState value and transforms it in a string
func getStateHelper(state libvirt.DomainState) string {
	ds := ""
	switch state {
	case libvirt.DOMAIN_NOSTATE:
		ds = "no state"
	case libvirt.DOMAIN_RUNNING:
		ds = "Running"
	case libvirt.DOMAIN_BLOCKED:
		ds = "Blocked"
	case libvirt.DOMAIN_CRASHED:
		ds = "Crashed"
	case libvirt.DOMAIN_SHUTDOWN:
	case libvirt.DOMAIN_SHUTOFF:
		ds = "Shutdown"
	case libvirt.DOMAIN_PMSUSPENDED:
		ds = "Suspended"
	case libvirt.DOMAIN_PAUSED:
		ds = "Paused"
	default:
		ds = "n/a"
	}
	return ds
}

// Some of it might not be needed anymore...
// This is where we get the interface name and its IP
func getInterfaceSpecs(dom libvirt.Domain, vmname string) (string, string, *ce.CustomError) {
	var domainInterface []libvirt.DomainInterface
	var interfaceName, interfaceAddress string
	var err error

	domainInterface, err = dom.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_AGENT)
	if err != nil {
		return "", "", &ce.CustomError{Title: "Unable to get interface specs", Message: err.Error()}
	}
	for _, di := range domainInterface {
		if len(di.Name) > 2 && (di.Name[:3] == "enp" || di.Name[:3] == "eth") {
			interfaceName = di.Name
			domainIPaddresses := di.Addrs
			for _, dipa := range domainIPaddresses {
				if dipa.Type == libvirt.IP_ADDR_TYPE_IPV4 {
					interfaceAddress = dipa.Addr
				}
			}

		}
	}
	return interfaceName, interfaceAddress, nil
}

// getUptime() : gets the active VM's uptime from the database
func getUptime(lastState string) string {
	//stateInDB := time.Unix(lastState, 0)
	a, _ := strconv.ParseInt(lastState, 10, 64)
	//deltaUnix := time.Now().Unix() - lastState
	return time.Unix(time.Now().Unix()-a, 0).Format("15:04:05")
}
