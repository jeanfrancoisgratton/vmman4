// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmSetResources.go
// Original timestamp : 2026.09.11 17:13:49

package vm_mgt

import (
	"fmt"
	"strconv"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// SetVMem : sets a VM's memory, in MiB. memParams holds 1 or 2 values: [min] or [min, max].
// When only one value is passed, min = max = that value.
// Overcommitting the hypervisor's physical memory is only warned about, not blocked.
func SetVMem(vmname string, memParams []string) *ce.CustomError {
	var conn *libvirt.Connect
	var domain *libvirt.Domain
	var err *ce.CustomError

	minMem, e := strconv.ParseUint(memParams[0], 10, 64)
	if e != nil {
		return &ce.CustomError{Title: "Invalid minimum memory value", Message: e.Error()}
	}
	maxMem := minMem
	if len(memParams) == 2 {
		if maxMem, e = strconv.ParseUint(memParams[1], 10, 64); e != nil {
			return &ce.CustomError{Title: "Invalid maximum memory value", Message: e.Error()}
		}
	}
	if minMem == 0 {
		return &ce.CustomError{Title: "Invalid memory value", Message: "Memory must be greater than 0"}
	}
	if minMem > maxMem {
		return &ce.CustomError{Title: "Invalid memory range", Message: "Minimum memory cannot exceed maximum memory"}
	}

	if err = connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = shared.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	if domain, err = shared.GetDomain(conn, vmname); err != nil {
		return err
	}
	defer domain.Free()

	Wait4Shutdown(domain, vmname)

	nodeInfo, nerr := conn.GetNodeInfo()
	if nerr != nil {
		return &ce.CustomError{Title: "Unable to fetch hypervisor node info", Message: nerr.Error()}
	}

	var warning *ce.CustomError
	if maxMem*1024 > nodeInfo.Memory {
		warning = &ce.CustomError{Fatality: ce.Warning, Title: "Memory overcommit on " + vmname,
			Message: fmt.Sprintf("requested max memory (%d MiB) exceeds the hypervisor's physical memory (%d MiB)", maxMem, nodeInfo.Memory/1024)}
	}

	if serr := domain.SetMemoryFlags(maxMem*1024, libvirt.DOMAIN_MEM_MAXIMUM|libvirt.DOMAIN_MEM_CONFIG); serr != nil {
		return &ce.CustomError{Title: "Cannot set max memory for " + vmname, Message: serr.Error()}
	}
	if serr := domain.SetMemoryFlags(minMem*1024, libvirt.DOMAIN_MEM_CONFIG); serr != nil {
		return &ce.CustomError{Title: "Cannot set memory for " + vmname, Message: serr.Error()}
	}

	fmt.Println(hftx.EnabledSign(vmname + hftx.Green(fmt.Sprintf(" memory set: min=%d MiB, max=%d MiB", minMem, maxMem))))

	if warning != nil {
		return warning
	}
	return nil
}

// SetVCPUs : sets a VM's vCPU count.
// Overcommitting the hypervisor's physical CPUs is only warned about, not blocked.
func SetVCPUs(vmname, count string) *ce.CustomError {
	var conn *libvirt.Connect
	var domain *libvirt.Domain
	var err *ce.CustomError

	vcpus, e := strconv.ParseUint(count, 10, 32)
	if e != nil {
		return &ce.CustomError{Title: "Invalid vCPU count", Message: e.Error()}
	}
	if vcpus == 0 {
		return &ce.CustomError{Title: "Invalid vCPU count", Message: "vCPU count must be greater than 0"}
	}

	if err = connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = shared.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	if domain, err = shared.GetDomain(conn, vmname); err != nil {
		return err
	}
	defer domain.Free()

	Wait4Shutdown(domain, vmname)

	nodeInfo, nerr := conn.GetNodeInfo()
	if nerr != nil {
		return &ce.CustomError{Title: "Unable to fetch hypervisor node info", Message: nerr.Error()}
	}

	var warning *ce.CustomError
	if uint(vcpus) > nodeInfo.Cpus {
		warning = &ce.CustomError{Fatality: ce.Warning, Title: "vCPU overcommit on " + vmname,
			Message: fmt.Sprintf("requested %d vCPUs exceeds the hypervisor's %d physical CPUs", vcpus, nodeInfo.Cpus)}
	}

	if serr := domain.SetVcpusFlags(uint(vcpus), libvirt.DOMAIN_VCPU_MAXIMUM|libvirt.DOMAIN_VCPU_CONFIG); serr != nil {
		return &ce.CustomError{Title: "Cannot set max vCPUs for " + vmname, Message: serr.Error()}
	}
	if serr := domain.SetVcpusFlags(uint(vcpus), libvirt.DOMAIN_VCPU_CONFIG); serr != nil {
		return &ce.CustomError{Title: "Cannot set vCPUs for " + vmname, Message: serr.Error()}
	}

	fmt.Println(hftx.EnabledSign(vmname + hftx.Green(fmt.Sprintf(" vCPUs set: %d", vcpus))))

	if warning != nil {
		return warning
	}
	return nil
}
