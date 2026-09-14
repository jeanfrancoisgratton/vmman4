// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/cluster_mgt/listClusters.go

package cluster_mgt

import (
	"encoding/json"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"github.com/jeanfrancoisgratton/helperFunctions/v5/prettyjson"
)

// ListClusters pretty-prints the contents of clusters.json.
func ListClusters() *ce.CustomError {
	cf, cerr := loadClustersFile()
	if cerr != nil {
		return cerr
	}

	jStream, err := json.Marshal(cf)
	if err != nil {
		return &ce.CustomError{Title: "Error marshaling JSON", Message: err.Error()}
	}

	if err := prettyjson.Print(jStream); err != nil {
		return &ce.CustomError{Title: "Error printing " + clustersFilePath(), Message: err.Error()}
	}
	return nil
}
