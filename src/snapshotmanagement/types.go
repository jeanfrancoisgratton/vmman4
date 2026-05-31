// vmman4
// Written by J.F. Gratton (jean-francois@famillegratton.net)
// Original filename : src/snapshotmanagement/types.go
// original timestamp : 2026/05/31 15:00:31

package snapshotmanagement

import "encoding/xml"

// Snapshot XML definitions
type ParentElement struct {
	XMLName    xml.Name `xml:"parent"`
	ParentName string   `xml:"name"`
}
type SnapshotXMLstruct struct {
	SnapshotName    string        `xml:"name"`
	CreationTime    int64         `xml:"creationTime"`
	Parent          ParentElement `xml:"parent"`
	CurrentSnapshot bool
}
