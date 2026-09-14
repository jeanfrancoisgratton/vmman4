// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/types_create.go
// Original timestamp : 2026.09.12 22:25:21

package vm_mgt

import "encoding/xml"

// ─── Types ────────────────────────────────────────────────────────────────────

// VMSpec is unmarshalled from the JSON file passed to Create.
type VMSpec struct {
	Name   string `json:"name"`
	Memory uint   `json:"memory_mb"` // MiB
	VCPUs  uint   `json:"vcpus"`

	// Arch and Machine are optional. When omitted, libvirt picks the best match
	// for the host. Specify them only for cross-arch VMs or a pinned machine
	// version, e.g. "aarch64"/"virt" or "x86_64"/"pc-q35-8.2".
	Arch    string `json:"arch,omitempty"`
	Machine string `json:"machine,omitempty"`

	Disks []struct {
		PoolName   string `json:"pool"`   // storage pool name
		VolumeName string `json:"volume"` // volume name inside that pool
		DeviceName string `json:"device"` // e.g. "vda"
		Bus        string `json:"bus"`    // e.g. "virtio", "ide", "scsi"

		// SizeGB is the volume's capacity, in GiB. It is only read -- and only
		// required -- when VolumeName does not already exist in PoolName: in
		// that case a new qcow2 volume of this size is created. If the volume
		// already exists, SizeGB is ignored and the existing volume is used as-is.
		SizeGB float64 `json:"size_gb,omitempty"`
	} `json:"disks"`

	Networks []struct {
		Source string `json:"source"` // network name, e.g. "default"
		Model  string `json:"model"`  // e.g. "virtio", "e1000"
	} `json:"networks"`

	// OSXML is an optional raw <os>…</os> XML snippet. When omitted, Create
	// builds one from the resolved arch and machine type.
	OSXML string `json:"os_xml,omitempty"`
}

// ─── Internal capability XML structs ─────────────────────────────────────────

type domainCapsXML struct {
	XMLName xml.Name `xml:"domainCapabilities"`
	Path    string   `xml:"path"`
	Arch    string   `xml:"arch"`
	Machine string   `xml:"machine"`
}
