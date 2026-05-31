// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vmmanagement/vmCreate.go

package vmmanagement

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

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

// ─── Public entry point ───────────────────────────────────────────────────────

// Create reads a VMSpec from a JSON file, resolves all host-specific parameters
// (emulator binary, arch, machine type) via libvirt's domain capabilities API,
// and defines the domain. The VM is NOT started; call dom.Create() to boot it.
//
// Minimal JSON spec — arch, machine, and os_xml are all optional:
//
//	{
//	  "name":      "my-vm",
//	  "memory_mb": 2048,
//	  "vcpus":     2,
//	  "disks": [
//	    { "pool": "default", "volume": "my-vm.qcow2", "device": "vda", "bus": "virtio" }
//	  ],
//	  "networks": [
//	    { "source": "default", "model": "virtio" }
//	  ]
//	}
func Create(conn *libvirt.Connect, specFile string) (*libvirt.Domain, *ce.CustomError) {
	spec, cerr := loadSpec(specFile)
	if cerr != nil {
		return nil, cerr
	}

	// Resolve the emulator path from libvirt capabilities — distro-agnostic.
	emulator, cerr := resolveEmulator(conn, spec.Arch, spec.Machine)
	if cerr != nil {
		return nil, cerr
	}

	// Fill in arch/machine defaults from the same capabilities document.
	resolvedArch, resolvedMachine, cerr := resolveArchMachine(conn, spec.Arch, spec.Machine)
	if cerr != nil {
		return nil, cerr
	}

	disksXML, cerr := buildDisksXML(conn, spec)
	if cerr != nil {
		return nil, cerr
	}

	nicsXML := buildNICsXML(spec)

	osXML := spec.OSXML
	if osXML == "" {
		osXML = fmt.Sprintf(
			`<os>
    <type arch="%s" machine="%s">hvm</type>
    <boot dev="hd"/>
  </os>`, resolvedArch, resolvedMachine)
	}

	domXMLStr := fmt.Sprintf(`<domain type="kvm">
  <name>%s</name>
  <memory unit="MiB">%d</memory>
  <currentMemory unit="MiB">%d</currentMemory>
  <vcpu placement="static">%d</vcpu>
  %s
  <features>
    <acpi/>
    <apic/>
  </features>
  <cpu mode="host-passthrough" check="none"/>
  <devices>
    <emulator>%s</emulator>
    %s
    %s
    <serial type="pty"><target port="0"/></serial>
    <console type="pty"><target type="serial" port="0"/></console>
    <channel type="unix">
      <target type="virtio" name="org.qemu.guest_agent.0"/>
    </channel>
    <input type="tablet" bus="usb"/>
    <graphics type="vnc" port="-1" autoport="yes"/>
    <video><model type="vga" vram="16384" heads="1"/></video>
    <memballoon model="virtio"/>
  </devices>
</domain>`,
		spec.Name,
		spec.Memory, spec.Memory,
		spec.VCPUs,
		osXML,
		emulator,
		disksXML,
		nicsXML,
	)

	dom, err := conn.DomainDefineXML(domXMLStr)
	if err != nil {
		return nil, &ce.CustomError{Title: "Create: DomainDefineXML failed", Message: err.Error()}
	}

	return dom, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// loadSpec reads and validates a VMSpec from a JSON file.
func loadSpec(specFile string) (*VMSpec, *ce.CustomError) {
	data, err := os.ReadFile(specFile)
	if err != nil {
		return nil, &ce.CustomError{Title: "loadSpec: cannot read file", Message: err.Error()}
	}

	var spec VMSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, &ce.CustomError{Title: "loadSpec: invalid JSON", Message: err.Error()}
	}
	if spec.Name == "" {
		return nil, &ce.CustomError{Title: "loadSpec: 'name' is required"}
	}
	if spec.Memory == 0 {
		return nil, &ce.CustomError{Title: "loadSpec: 'memory_mb' must be > 0"}
	}
	if spec.VCPUs == 0 {
		spec.VCPUs = 1
	}

	return &spec, nil
}

// resolveEmulator queries libvirt's domain capabilities for this host and
// returns the correct emulator binary path. This is distro-agnostic: Debian
// uses /usr/bin/qemu-system-x86_64, RHEL/Rocky use /usr/libexec/qemu-kvm,
// and so on — libvirt always knows the right path regardless of distro.
// Falls back from kvm to qemu virttype automatically when KVM is unavailable.
func resolveEmulator(conn *libvirt.Connect, arch, machine string) (string, *ce.CustomError) {
	caps, cerr := getDomainCaps(conn, arch, machine)
	if cerr != nil {
		return "", cerr
	}
	if caps.Path == "" {
		return "", &ce.CustomError{Title: "resolveEmulator: capabilities XML contains no emulator path"}
	}
	return caps.Path, nil
}

// resolveArchMachine returns the effective arch and machine type, filling in
// libvirt's host defaults when the caller left either field empty.
func resolveArchMachine(conn *libvirt.Connect, arch, machine string) (string, string, *ce.CustomError) {
	if arch != "" && machine != "" {
		return arch, machine, nil
	}
	caps, cerr := getDomainCaps(conn, arch, machine)
	if cerr != nil {
		return "", "", cerr
	}
	if arch == "" {
		arch = caps.Arch
	}
	if machine == "" {
		machine = caps.Machine
	}
	if arch == "" || machine == "" {
		return "", "", &ce.CustomError{
			Title:   "resolveArchMachine: libvirt returned no arch/machine in capabilities",
			Message: fmt.Sprintf("got arch=%q machine=%q", arch, machine),
		}
	}
	return arch, machine, nil
}

// getDomainCaps is the single point of contact with GetDomainCapabilities.
// It tries kvm first, then falls back to qemu (e.g. nested virt / CI).
func getDomainCaps(conn *libvirt.Connect, arch, machine string) (*domainCapsXML, *ce.CustomError) {
	capsXML, err := conn.GetDomainCapabilities("", arch, machine, "kvm", 0)
	if err != nil {
		capsXML, err = conn.GetDomainCapabilities("", arch, machine, "qemu", 0)
		if err != nil {
			return nil, &ce.CustomError{Title: "getDomainCaps: GetDomainCapabilities failed", Message: err.Error()}
		}
	}
	var caps domainCapsXML
	if err := xml.Unmarshal([]byte(capsXML), &caps); err != nil {
		return nil, &ce.CustomError{Title: "getDomainCaps: cannot parse capabilities XML", Message: err.Error()}
	}
	return &caps, nil
}

// buildDisksXML resolves each disk's volume path from its pool and assembles
// the <disk> XML elements.
func buildDisksXML(conn *libvirt.Connect, spec *VMSpec) (string, *ce.CustomError) {
	var out string

	for i, d := range spec.Disks {
		if d.PoolName == "" || d.VolumeName == "" {
			return "", &ce.CustomError{
				Title:   fmt.Sprintf("buildDisksXML: disk[%d] missing pool or volume", i),
				Message: "'pool' and 'volume' are required for every disk entry",
			}
		}

		pool, err := conn.LookupStoragePoolByName(d.PoolName)
		if err != nil {
			return "", &ce.CustomError{
				Title:   fmt.Sprintf("buildDisksXML: pool %q not found", d.PoolName),
				Message: err.Error(),
			}
		}

		vol, err := pool.LookupStorageVolByName(d.VolumeName)
		pool.Free()
		if err != nil {
			return "", &ce.CustomError{
				Title:   fmt.Sprintf("buildDisksXML: volume %q not found in pool %q", d.VolumeName, d.PoolName),
				Message: err.Error(),
			}
		}

		volPath, err := vol.GetPath()
		vol.Free()
		if err != nil {
			return "", &ce.CustomError{
				Title:   fmt.Sprintf("buildDisksXML: cannot get path for volume %q", d.VolumeName),
				Message: err.Error(),
			}
		}

		bus := d.Bus
		if bus == "" {
			bus = "virtio"
		}
		dev := d.DeviceName
		if dev == "" {
			dev = fmt.Sprintf("vd%c", 'a'+i)
		}

		out += fmt.Sprintf(`
    <disk type="file" device="disk">
      <driver name="qemu" type="qcow2" cache="none" io="native"/>
      <source file="%s"/>
      <target dev="%s" bus="%s"/>
    </disk>`, volPath, dev, bus)
	}

	return out, nil
}

// buildNICsXML assembles the <interface> XML elements.
func buildNICsXML(spec *VMSpec) string {
	var out string
	for _, n := range spec.Networks {
		model := n.Model
		if model == "" {
			model = "virtio"
		}
		out += fmt.Sprintf(`
    <interface type="network">
      <source network="%s"/>
      <model type="%s"/>
    </interface>`, n.Source, model)
	}
	return out
}
