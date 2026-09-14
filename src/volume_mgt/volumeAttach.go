// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/volumeAttach.go

package volume_mgt

import (
	"encoding/xml"
	"fmt"

	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// AttachVolume attaches an existing volume (poolName/volName) to vmName as a
// new virtio disk, on the next free vd* target device. Since the disk is
// attached to the persistent (inactive) domain config, the VM is shut down
// first -- gracefully, then forcefully after 15s -- if it's running.
func AttachVolume(vmName, poolName, volName string) *ce.CustomError {
	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	domain, err := shared.GetDomain(conn, vmName)
	if err != nil {
		return err
	}
	defer domain.Free()

	pool, perr := conn.LookupStoragePoolByName(poolName)
	if perr != nil {
		return &ce.CustomError{Title: fmt.Sprintf("AttachVolume: pool %q not found", poolName), Message: perr.Error()}
	}
	defer pool.Free()

	vol, verr := pool.LookupStorageVolByName(volName)
	if verr != nil {
		return &ce.CustomError{Title: fmt.Sprintf("AttachVolume: volume %q not found in pool %q", volName, poolName), Message: verr.Error()}
	}
	defer vol.Free()

	volPath, verr := vol.GetPath()
	if verr != nil {
		return &ce.CustomError{Title: fmt.Sprintf("AttachVolume: cannot get path for volume %q", volName), Message: verr.Error()}
	}

	shared.Wait4Shutdown(domain, vmName)

	dev, cerr := nextFreeDiskTarget(domain)
	if cerr != nil {
		return cerr
	}

	diskXML := fmt.Sprintf(`<disk type="file" device="disk">
  <driver name="qemu" type="qcow2" cache="none" io="native"/>
  <source file="%s"/>
  <target dev="%s" bus="virtio"/>
</disk>`, volPath, dev)

	if aerr := domain.AttachDeviceFlags(diskXML, libvirt.DOMAIN_DEVICE_MODIFY_CONFIG); aerr != nil {
		return &ce.CustomError{
			Title:   fmt.Sprintf("AttachVolume: cannot attach %q to %q", volName, vmName),
			Message: aerr.Error(),
		}
	}

	fmt.Println(hftx.EnabledSign(fmt.Sprintf("%s/%s", poolName, volName) + hftx.Green(" ATTACHED") + fmt.Sprintf(" to %s as %s", vmName, dev)))
	return nil
}

// nextFreeDiskTarget parses the domain's (inactive) XML and returns the first
// unused "vd?" target device, matching the naming convention vm_mgt.buildDisksXML
// uses when a VM is first created.
func nextFreeDiskTarget(domain *libvirt.Domain) (string, *ce.CustomError) {
	xmlDesc, err := domain.GetXMLDesc(libvirt.DOMAIN_XML_INACTIVE)
	if err != nil {
		return "", &ce.CustomError{Title: "nextFreeDiskTarget: cannot get domain XML", Message: err.Error()}
	}

	var dXML domainXML
	if err := xml.Unmarshal([]byte(xmlDesc), &dXML); err != nil {
		return "", &ce.CustomError{Title: "nextFreeDiskTarget: cannot parse domain XML", Message: err.Error()}
	}

	used := make(map[string]bool, len(dXML.Devices.Disks))
	for _, d := range dXML.Devices.Disks {
		used[d.Target.Dev] = true
	}

	for c := 'a'; c <= 'z'; c++ {
		dev := fmt.Sprintf("vd%c", c)
		if !used[dev] {
			return dev, nil
		}
	}
	return "", &ce.CustomError{Title: "nextFreeDiskTarget: no free target device", Message: "domain already has 26 virtio disks (vda..vdz)"}
}
