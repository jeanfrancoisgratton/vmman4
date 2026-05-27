// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/connection/uriHandler.go
// Original timestamp: 2026/05/19 19:49:39

package connection

import (
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"vmman4/shared"
)

// ResolveConnectionURI : this is where we set the connection string up
// The connection info can be provided with -c CONNECTION_NAME or -C CONNECTION_STRING
func ResolveConnectionURI() *ce.CustomError {
	// First check, we need to check the value of the -c flag
	if shared.ConnectionFilename != "" {
		var ct ConnectionType
		if err := ct.LoadConnectionInfo(); err != nil {
			return &ce.CustomError{Title: "Unable to setup the connection information", Message: err.Error()}
		}

		if strings.ToLower(ct.Name) == "localhost" {
			shared.ConnectURI = "qemu:///system"
		} else {
			if ct.User == "" {
				return &ce.CustomError{Title: "Unable to build the connection string",
					Message: "Username is empty"}
			}
			shared.ConnectURI = "qemu+ssh://" + ct.User + "@" + ct.Host + "/system"
		}
	}
	return nil
}
