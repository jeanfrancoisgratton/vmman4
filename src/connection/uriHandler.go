// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/connection/uriHandler.go
// Original timestamp: 2026/05/19 19:49:39

package connection

import (
	"errors"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
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
			if ct.User != "" {
				return &ce.CustomError{Title: "Unable to build the connection string",
					Message: "Invalid username: " + ct.User}
			}
			shared.ConnectURI = "qemu+ssh://" + ct.User + "@" + ct.Host
		}
	}
	return nil
}

// Connect2HVM : this is where we actually handle the connection
func Connect2HVM() (*libvirt.Connect, *ce.CustomError) {
	conn, err := libvirt.NewConnect(shared.ConnectURI)
	if err != nil {
		var lverr libvirt.Error
		ok := errors.As(err, &lverr)
		if ok && (lverr.Message == "End of file while reading data: virt-ssh-helper: cannot connect to '/var/run/libvirt/libvirt-sock': Failed to connect socket to '/var/run/libvirt/libvirt-sock': Connection refused: Input/output error" ||
			strings.HasPrefix(lverr.Message, "Cannot recv data: ssh: connect to host")) {
			return nil, &ce.CustomError{Title: "Unable to connect to host", Message: "Hypervisor " + shared.ConnectURI + " is offline"}
		} else {
			return nil, &ce.CustomError{Title: "Unable to connect to host", Message: err.Error()}
		}
	}
	return conn, nil
}
