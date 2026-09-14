// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmHelpers.go
// Original timestamp: 2026/05/26 19:08:41

package vm_mgt

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"vmman4/shared"
	"vmman4/snapshot_mgt"
	"vmman4/volume_mgt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"golang.org/x/term"
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
		if isVirtualInterface(di.Name) {
			continue
		}
		for _, dipa := range di.Addrs {
			if dipa.Type == libvirt.IP_ADDR_TYPE_IPV4 {
				interfaceName = di.Name
				interfaceAddress = dipa.Addr
				break
			}
		}
		if interfaceName != "" {
			break
		}
	}
	return interfaceName, interfaceAddress, nil
}

// isVirtualInterface : filters out loopback and container/bridge interfaces
// (e.g. docker0, veth*) that guest agents report alongside the real NIC, so
// they don't get mistaken for the VM's actual network interface.
func isVirtualInterface(name string) bool {
	if len(name) < 2 {
		return true
	}
	virtualPrefixes := []string{"lo", "docker", "veth", "br-", "virbr", "vnet", "tap", "tun", "cni", "flannel", "cali", "podman"}
	for _, p := range virtualPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
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

	if doms, serr = listDomains(conn); serr != nil {
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
			if i.viCurrentSnapshot, cerr = snapshot_mgt.GetCurrentSnapshotName(conn, i.viName); cerr != nil {
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

		// DISK INFO
		storageInfo, serr := volume_mgt.GetStorageSpecs4VM(i.viName, conn)
		if serr != nil {
			return nil, serr
		}
		var diskTotalSize uint64
		for _, disk := range storageInfo.Disks {
			diskTotalSize += disk.SizeBytes
		}

		vmspec = append(vmspec, vmInfo{viId: i.viId, viName: i.viName, viState: getStateHelper(dState), viMem: specs.Memory / 1024, viCpu: specs.NrVirtCpu,
			viSnapshots: uint(numsnap), viCurrentSnapshot: i.viCurrentSnapshot, viInterfaceName: i.viInterfaceName, viIPaddress: i.viIPaddress,
			viDiskCount: uint(storageInfo.DiskCount), viDiskTotalSize: diskTotalSize})
	}
	return vmspec, nil
}

// termWidth returns the current terminal width, falling back to 80 on error.
func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
