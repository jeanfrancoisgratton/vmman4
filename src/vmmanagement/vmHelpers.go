// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vmmanagement/vmHelpers.go
// Original timestamp: 2026/05/26 19:08:41

package vmmanagement

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
	"vmman4/shared"
	"vmman4/snapshotmanagement"
)

// Wait4Shutdown : Tries 15 seconds to gracefully shutdown the VM, if not it will shutdown forcefully
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

// This is where we get the interface name and its IP
func getInterfaceSpecs(dom libvirt.Domain, vmname string) (string, string, *ce.CustomError) {
	var domainInterface []libvirt.DomainInterface
	var interfaceName, interfaceAddress string
	var err error

	domainInterface, err = dom.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_AGENT)
	if err != nil {
		em := err.Error()
		if strings.Contains(em, "QEMU guest agent is not connected") {
			fmt.Println(hftx.WarningSign(vmname + " : QEMU guest agent is not connected"))
		} else {
			return "", "", &ce.CustomError{Title: "Unable to get interface specs", Message: err.Error()}
		}
	}
	for _, di := range domainInterface {
		//if len(di.Name) > 2 && (di.Name[:3] == "enp" || di.Name[:3] == "eth") {
		if len(di.Name) > 2 {
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

// getUptime : gets the active VM's uptime from the database
func getUptime(lastState string) string {
	//stateInDB := time.Unix(lastState, 0)
	a, _ := strconv.ParseInt(lastState, 10, 64)
	//deltaUnix := time.Now().Unix() - lastState
	return time.Unix(time.Now().Unix()-a, 0).Format("15:04:05")
}

// collectInfo: Lists all VMs and collects info on each of them
func collectInfo(conn *libvirt.Connect) ([]vmInfo, *ce.CustomError) {
	var (
		snapshotflags libvirt.DomainSnapshotListFlags
		numsnap       int
		vmspec        []vmInfo
		i             vmInfo
		dState        libvirt.DomainState
		doms          []libvirt.Domain
		cerr          *ce.CustomError
		serr          *ce.CustomError
		domain        *libvirt.Domain
	)

	if doms, serr = GetVMlist(); serr != nil {
		return nil, serr
	}

	for _, dom := range doms {
		specs, err := dom.GetInfo()
		i.viId, err = dom.GetID()
		if err != nil {
			i.viId = 0
		}
		// VM NAME
		i.viName, _ = dom.GetName()
		// VM STATE
		dState, _, _ = dom.GetState()
		i.viState = getStateHelper(dState)
		if i.viState == "Running" {
			// INTERFACE INFO
			i.viInterfaceName, i.viIPaddress, cerr = getInterfaceSpecs(dom, i.viName)
			if cerr != nil {
				return nil, cerr
			}
		} else {
			i.viInterfaceName = ""
			i.viIPaddress = ""
		}
		// SNAPSHOT INFO
		if domain, cerr = shared.GetDomain(conn, i.viName); cerr != nil {
			return nil, cerr
		}
		if domain == nil {
			continue
		}
		defer domain.Free()
		numsnap, _ = domain.SnapshotNum(snapshotflags)
		if numsnap > 0 {
			if i.viCurrentSnapshot, cerr = snapshotmanagement.GetCurrentSnapshotName(conn, i.viName); cerr != nil {
				return nil, cerr
			}
		} else {
			i.viCurrentSnapshot = "n/a"
		}

		// uptime is not yet working, so commenting out that block
		//if i.viId > 0 {
		//	// time.Unix(time.Now().Unix() - lastState, 0).Format("2006.01.02 15:04:05")
		//	i.viLastStatusChange = getUptime(i.viLastStatusChange)
		//}

		vmspec = append(vmspec, vmInfo{viId: i.viId, viName: i.viName, viState: getStateHelper(dState), viMem: specs.Memory / 1024, viCpu: specs.NrVirtCpu,
			viSnapshots: uint(numsnap), viCurrentSnapshot: i.viCurrentSnapshot, viInterfaceName: i.viInterfaceName, viIPaddress: i.viIPaddress})
	}
	return vmspec, nil
}
