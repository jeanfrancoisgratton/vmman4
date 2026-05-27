// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/connection/types.go
// Original timestamp: 2026/05/17 22:21:02

package connection

var (
	ConnectionHost    string
	ConnectionUser    string
	ConnectionComment string
)

// The ConnectionType struct holds all the relevant info on how to connect on a specific hypervisor
type ConnectionType struct {
	Name     string `json:"connectionName"`
	Host     string `json:"connectionHost,omitempty"`
	User     string `json:"connectionUser"`
	Comments string `json:"connectionComments,omitempty"`
}
