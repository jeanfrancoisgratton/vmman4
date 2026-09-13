// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/connection_mgt/createRemoveConnections.go
// Original timestamp: 2026/05/19 19:54:47

package connection_mgt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// CreateConnection : Create a connection_mgt file; user will be prompted for relevant information,
// unless the --name, --host --user flags are used
func CreateConnection() *ce.CustomError {
	ct := ConnectionType{}

	fmt.Println("Please enter a connection_mgt name")
	ct.Name = hf.GetStringValFromPrompt("This will be the filename of the connection_mgt definition file, located in $HOME/.config/JFG/vmman4/ : ")

	if ConnectionHost == "" {
		ct.Host = hf.GetStringValFromPrompt("Please enter a host name (leave blank for localhost): ")
	}
	if ConnectionUser == "" {
		ct.User = hf.GetStringValFromPrompt("Please enter a username: ")
	}
	if ct.Comments == "" {
		ct.Comments = hf.GetStringValFromPrompt("[OPTIONAL] Please enter a comment: ")
	}

	return ct.SaveConnection()
}

// RemoveConnection : Remove one or many connection_mgt file(s)
func RemoveConnection(connfiles []string) *ce.CustomError {
	for _, connfile := range connfiles {
		if !strings.HasSuffix(connfile, ".json") {
			connfile += ".json"
		}
		if err := os.Remove(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", connfile)); err != nil {
			return &ce.CustomError{Title: "Error removing " + connfile, Message: err.Error()}
		}
		if !shared.QuietOutput {
			fmt.Println(hftx.EnabledSign("Removed " + connfile))
		}
	}
	return nil
}
