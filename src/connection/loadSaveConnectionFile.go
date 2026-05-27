// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/connection/loadSaveConnectionFile.go
// Original timestamp: 2026/05/19 20:05:24

package connection

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"vmman4/shared"
)

func (ct ConnectionType) ConnSave() *ce.CustomError {

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

// LoadConnectionInfo : the connection JSON file is unmarshalled and loaded in a variable
func (ct ConnectionType) LoadConnectionInfo() *ce.CustomError {
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
