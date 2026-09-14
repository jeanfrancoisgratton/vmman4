// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmInfo.go
// Original timestamp: 2026/09/11

package vm_mgt

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"vmman4/connection_mgt"
	"vmman4/shared"
	"vmman4/snapshot_mgt"
	"vmman4/volume_mgt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	"libvirt.org/go/libvirt"
)

// VmInfo : Displays detailed, single-VM information as tab-separated key/value pairs.
func VmInfo(vmname string) *ce.CustomError {
	var (
		conn *libvirt.Connect
		err  *ce.CustomError
	)

	if err = connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = shared.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	dom, err := shared.GetDomain(conn, vmname)
	if err != nil {
		return err
	}
	if dom == nil {
		return &ce.CustomError{Title: "VmInfo", Message: "VM " + vmname + " not found"}
	}
	defer dom.Free()

	domID, gerr := dom.GetID()
	if gerr != nil {
		domID = 0
	}

	dState, _, _ := dom.GetState()
	state := getStateHelper(dState)

	specs, gerr := dom.GetInfo()
	if gerr != nil {
		return &ce.CustomError{Title: "VmInfo: unable to get domain info", Message: gerr.Error()}
	}

	uuid, _ := dom.GetUUIDString()
	autostart, _ := dom.GetAutostart()
	persistent, _ := dom.IsPersistent()

	numsnap, _ := dom.SnapshotNum(0)
	currentSnapshot := "n/a"
	if numsnap > 0 {
		if currentSnapshot, err = snapshot_mgt.GetCurrentSnapshotName(conn, vmname); err != nil {
			return err
		}
	}

	var ifName, ipAddr string
	if state == "Running" {
		if ifName, ipAddr, err = getInterfaceSpecs(*dom, vmname); err != nil {
			return err
		}
	}

	storageInfo, cerr := volume_mgt.GetStorageSpecs4VM(vmname, conn)
	if cerr != nil {
		return cerr
	}

	var pools []string
	seenPool := map[string]bool{}
	var totalDiskBytes uint64
	for _, d := range storageInfo.Disks {
		totalDiskBytes += d.SizeBytes
		if d.PoolName != "" && d.PoolName != "n/a" && !seenPool[d.PoolName] {
			seenPool[d.PoolName] = true
			pools = append(pools, d.PoolName)
		}
	}
	poolDisplay := "n/a"
	if len(pools) > 0 {
		poolDisplay = strings.Join(pools, ", ")
	}

	minMem, verr := hf.BytesToUnit(specs.Memory*1024, 'g', 2)
	if verr != nil {
		return &ce.CustomError{Title: "VmInfo: unable to convert memory units", Message: verr.Error()}
	}
	maxMem, verr := hf.BytesToUnit(specs.MaxMem*1024, 'g', 2)
	if verr != nil {
		return &ce.CustomError{Title: "VmInfo: unable to convert memory units", Message: verr.Error()}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(w, "Hypervisor:\t%s\n", shared.ConnectURI)
	fmt.Fprintf(w, "Domain ID:\t%d\n", domID)
	fmt.Fprintf(w, "VM name:\t%s\n", vmname)
	fmt.Fprintf(w, "UUID:\t%s\n", uuid)
	fmt.Fprintf(w, "State:\t%s\n", state)
	fmt.Fprintf(w, "Persistent:\t%t\n", persistent)
	fmt.Fprintf(w, "Autostart:\t%t\n", autostart)
	fmt.Fprintf(w, "Min memory (GB):\t%s\n", minMem)
	fmt.Fprintf(w, "Max memory (GB):\t%s\n", maxMem)
	fmt.Fprintf(w, "vCPUs:\t%d\n", specs.NrVirtCpu)
	fmt.Fprintf(w, "Snapshots:\t%d\n", numsnap)
	fmt.Fprintf(w, "Current snapshot:\t%s\n", currentSnapshot)
	fmt.Fprintf(w, "Storage pool:\t%s\n", poolDisplay)
	fmt.Fprintf(w, "# disks:\t%d\n", storageInfo.DiskCount)
	for _, d := range storageInfo.Disks {
		sz, serr := hf.BytesToUnit(d.SizeBytes, 'g', 2)
		if serr != nil {
			return &ce.CustomError{Title: "VmInfo: unable to convert disk size units", Message: serr.Error()}
		}
		fmt.Fprintf(w, "  disk %s (%s):\t%s GB\n", d.Device, d.Type, sz)
	}
	totalSize, serr := hf.BytesToUnit(totalDiskBytes, 'g', 2)
	if serr != nil {
		return &ce.CustomError{Title: "VmInfo: unable to convert disk size units", Message: serr.Error()}
	}
	fmt.Fprintf(w, "  TOTAL:\t%s GB\n", totalSize)
	fmt.Fprintf(w, "Interface:\t%s\n", valueOrNA(ifName))
	fmt.Fprintf(w, "IP address:\t%s\n", valueOrNA(ipAddr))
	w.Flush()

	return nil
}

// valueOrNA : returns "n/a" for empty strings, otherwise the string itself.
func valueOrNA(s string) string {
	if s == "" {
		return "n/a"
	}
	return s
}
