// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/connection_mgt/loadSaveConnectionFile.go
// Original timestamp: 2026/05/19 20:05:24

package connection_mgt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

func (ct *ConnectionType) SaveConnection() *ce.CustomError {
	jStream, err := json.MarshalIndent(ct, "", "  ")
	if err != nil {
		return &ce.CustomError{Title: "Error marshaling JSON", Message: err.Error()}
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", ct.Name+".json")

	if err := os.WriteFile(rcFile, jStream, 0600); err != nil {
		return &ce.CustomError{Title: "Error writing file", Message: err.Error()}
	}

	return nil
}

// LoadConnectionInfo : the connection_mgt JSON file is unmarshalled and loaded in a variable
func (ct *ConnectionType) LoadConnectionInfo() *ce.CustomError {
	if !strings.HasSuffix(shared.ConnectionFilename, ".json") {
		shared.ConnectionFilename += ".json"
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", shared.ConnectionFilename)
	rc2json, err := os.ReadFile(rcFile)
	if err != nil {
		return &ce.CustomError{Title: "Error reading file", Message: err.Error()}
	}
	err = json.Unmarshal(rc2json, &ct)
	if err != nil {
		return &ce.CustomError{Title: "Error unmarshalling file", Message: err.Error()}
	}

	return nil
}
