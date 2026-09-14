// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/types.go

package volume_mgt

import "encoding/xml"

// DiskInfo describes a single disk attached to a VM.
type DiskInfo struct {
	Device         string // e.g. "vda", "sda"
	Type           string // e.g. "file", "block", "network"
	SourcePath     string // file path, pool volume, or network source
	SizeBytes      uint64 // logical capacity in bytes (0 if unavailable)
	PoolName       string // owning storage pool name, or "n/a" if none
	PoolTargetPath string // owning storage pool's target path, or "n/a" if none
	PoolState      string // owning storage pool's state, or "n/a" if none
}

// VMStorageInfo is the result of GetStorageSpecs4VM.
type VMStorageInfo struct {
	DomainName string
	DiskCount  int
	Disks      []DiskInfo
}

// VolumeInfo describes a volume inside a storage pool.
type VolumeInfo struct {
	Name      string
	Path      string
	SizeBytes uint64
}

// ─── Internal XML structs ─────────────────────────────────────────────────────

type domainXML struct {
	XMLName xml.Name   `xml:"domain"`
	Devices devicesXML `xml:"devices"`
}

type devicesXML struct {
	Disks []diskXML `xml:"disk"`
}

type diskXML struct {
	Type   string     `xml:"type,attr"`
	Device string     `xml:"device,attr"`
	Source diskSource `xml:"source"`
	Target diskTarget `xml:"target"`
}

type diskSource struct {
	File   string `xml:"file,attr"`
	Dev    string `xml:"dev,attr"`
	Name   string `xml:"name,attr"`
	Pool   string `xml:"pool,attr"`
	Volume string `xml:"volume,attr"`
}

type diskTarget struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}
