// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/loadSaveClustersFile.go

package cluster_mgt

import (
	"encoding/json"
	"os"
	"path/filepath"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// clustersFilePath is always $HOME/.config/JFG/vmman4/clusters.json, regardless
// of which hypervisor(s) the clusters it describes actually live on.
func clustersFilePath() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", "clusters.json")
}

// loadClustersFile reads clusters.json, returning an empty ClustersFile if it
// doesn't exist yet (first use).
func loadClustersFile() (ClustersFile, *ce.CustomError) {
	path := clustersFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ClustersFile{}, nil
		}
		return nil, &ce.CustomError{Title: "Error reading " + path, Message: err.Error()}
	}
	if len(data) == 0 {
		return ClustersFile{}, nil
	}

	cf := ClustersFile{}
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, &ce.CustomError{Title: "Error unmarshalling " + path, Message: err.Error()}
	}
	return cf, nil
}

func saveClustersFile(cf ClustersFile) *ce.CustomError {
	path := clustersFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return &ce.CustomError{Title: "Error creating config directory", Message: err.Error()}
	}

	jStream, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return &ce.CustomError{Title: "Error marshaling JSON", Message: err.Error()}
	}
	if err := os.WriteFile(path, jStream, 0600); err != nil {
		return &ce.CustomError{Title: "Error writing " + path, Message: err.Error()}
	}
	return nil
}
