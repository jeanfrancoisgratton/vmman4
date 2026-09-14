// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmCreateSample.go

package vm_mgt

import (
	"fmt"
	"os"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// sampleSpec is an annotated example of the JSON file `vm create` expects.
// It is NOT valid JSON (the // comments below are illegal in real JSON) --
// it exists purely to document every field for someone who has never
// touched libvirt before. Strip the comments and unneeded lines before
// using it as a real spec file.
const sampleSpec = `{
  // ── Required fields ──────────────────────────────────────────────────────

  // Name of the VM as it will appear in 'virsh list' / 'vm list'. Must be
  // unique on the target hypervisor.
  "name": "my-vm",

  // RAM to allocate, in MiB (1024 = 1 GiB). 2048 MiB = 2 GiB, 4096 MiB = 4 GiB.
  "memory_mb": 2048,

  // Number of virtual CPUs. Omit or set to 0 to default to 1.
  "vcpus": 2,

  // ── Optional fields ──────────────────────────────────────────────────────

  // CPU architecture. Leave empty ("") to let libvirt pick the host's native
  // arch. Only set this for cross-arch VMs, e.g. "aarch64", "x86_64".
  "arch": "",

  // QEMU machine type, e.g. "pc-q35-8.2", "virt" (aarch64). Leave empty to
  // let libvirt pick the best/newest match for the host. Pin this only if
  // you need a stable machine type across host upgrades (e.g. for migration).
  "machine": "",

  // ── Disks ─────────────────────────────────────────────────────────────────
  // A VM can have any number of disks -- just add more objects to this array,
  // one per drive. Each one is handled independently:
  //   * If "volume" already exists inside "pool" (e.g. created ahead of time
  //     with 'virsh vol-create-as', or by whatever provisioned the pool), it
  //     is attached as-is and "size_gb" is ignored.
  //   * If "volume" does NOT exist yet, vmman4 creates it for you as a new
  //     qcow2 volume -- and "size_gb" then becomes REQUIRED. Omitting it is
  //     an error when the volume doesn't already exist.
  "disks": [
    {
      // Name of an existing libvirt storage pool. See 'vmman pool list'.
      "pool": "default",

      // Volume (e.g. a .qcow2 file) inside that pool. This example assumes
      // "my-vm.qcow2" already exists in "default" -- it's attached as-is,
      // so "size_gb" is irrelevant here and can be omitted.
      "volume": "my-vm.qcow2",

      // Guest-visible device name. Follow libvirt's vd* convention for
      // virtio disks: vda, vdb, vdc, ... Optional -- if omitted, vmman4
      // assigns vda, vdb, vdc... in array order.
      "device": "vda",

      // Disk bus type. "virtio" is fastest and recommended for Linux guests
      // with virtio drivers. Use "sata" or "ide" for older/legacy guests
      // that lack virtio drivers (e.g. plain install media). Optional,
      // defaults to "virtio".
      "bus": "virtio"
    },
    {
      // Second example: a data disk that does NOT exist yet in the pool.
      // Because "my-vm-data.qcow2" isn't found in "default", vmman4 will
      // create it as a new 50 GiB qcow2 volume before attaching it.
      "pool": "default",
      "volume": "my-vm-data.qcow2",
      "device": "vdb",
      "bus": "virtio",

      // Capacity of the new volume, in GiB. REQUIRED because the volume
      // above doesn't exist yet. Ignored (and safe to omit) for disks that
      // reference an already-existing volume, like the first one above.
      "size_gb": 50
    }
    // Add as many more disk objects here as the VM needs.
  ],

  // ── Networks ──────────────────────────────────────────────────────────────
  "networks": [
    {
      // Name of an existing libvirt network (see 'virsh net-list'). "default"
      // is libvirt's built-in NATed network, present on most hosts out of
      // the box.
      "source": "default",

      // Virtual NIC model presented to the guest. "virtio" is fastest and
      // needs virtio-net drivers in the guest. Use "e1000" for guests
      // without virtio drivers. Optional, defaults to "virtio".
      "model": "virtio"
    }
    // Add more network objects here to give the VM multiple NICs.
  ],

  // ── Advanced / escape hatch ─────────────────────────────────────────────
  // Raw <os>...</os> libvirt domain XML, used verbatim instead of the
  // auto-generated one if provided. Only needed for unusual boot setups
  // (e.g. UEFI/OVMF, custom boot order, netboot). Leave empty ("") for a
  // normal BIOS boot using the resolved arch/machine above:
  //
  //   <os>
  //     <type arch="x86_64" machine="pc-q35-8.2">hvm</type>
  //     <boot dev="hd"/>
  //   </os>
  //
  // Example UEFI override:
  //   "os_xml": "<os><type arch=\"x86_64\" machine=\"q35\">hvm</type><loader readonly=\"yes\" type=\"pflash\">/usr/share/OVMF/OVMF_CODE.fd</loader><boot dev=\"hd\"/></os>"
  "os_xml": ""
}
`

// WriteSample writes the annotated example spec above to destPath, so a user
// unfamiliar with libvirt has a fully-documented starting point. The file it
// produces is illustrative, not valid JSON -- comments must be stripped and
// unwanted fields removed before it can be passed to 'vm create'.
func WriteSample(destPath string) *ce.CustomError {
	if err := os.WriteFile(destPath, []byte(sampleSpec), 0644); err != nil {
		return &ce.CustomError{Title: "WriteSample: cannot write file", Message: err.Error()}
	}
	fmt.Println(hftx.EnabledSign("Sample spec written to " + hftx.Blue(destPath)))
	fmt.Println(hftx.WarningSign("This file has // comments and is NOT valid JSON -- edit it into a real spec before using it with 'vm create'"))
	return nil
}
