// vmman3 : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/inventory/collectInfo.go
// 2022-09-17 14:01:35

package inventory

import (
	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
	"vmman4/shared"
	"vmman4/snapshotmanagement"
)

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

	if doms, serr = shared.GetVMlist(); serr != nil {
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
