// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename : structs.go
// Original timestamp : 2024/07/29 09:46

package db

// table: hypervisors
type DbHypervisors struct {
	ID             uint8  `json:"id" yaml:"id"`
	Name           string `json:"name" yaml:"name"`
	Address        string `json:"address" yaml:"address"`
	Connectinguser string `json:"connectinguser" yaml:"connectinguser"`
}

// table: storagepools
type DbStoragePools struct {
	ID    uint8  `json:"id" yaml:"id"`
	Name  string `json:"name" yaml:"name"`
	Path  string `json:"path" yaml:"path"`
	Owner string `json:"owner" yaml:"owner"`
}

// table: vmstates
type dbVmStates struct {
	ID              uint8  `json:"id" yaml:"id"`
	Name            string `json:"name" yaml:"name"`
	IP              string `json:"ip" yaml:"ip"`
	Online          bool   `json:"online" yaml:"online"`
	LastStateChange string `json:"laststatechange" yaml:"laststatechange"`
	OperatingSystem string `json:"os" yaml:"vmos"`
	Hypervisor      string `json:"hypervisor" yaml:"hypervisor"`
	StoragePool     string `json:"storagepool" yaml:"storagepool"`
}

// table: templates
type dbTemplates struct {
	ID              uint8  `json:"id" yaml:"id"`
	Name            string `json:"name" yaml:"name"`
	Owner           string `json:"owner" yaml:"owner"`
	StoragePool     string `json:"storagepool" yaml:"storagepool"`
	OperatingSystem string `json:"os" yaml:"os"`
}

type dbDisks struct {
	ID          uint   `json:"id" yaml:"id"`
	Name        string `json:"name" yaml:"name"`
	StoragePool string `json:"storagepool" yaml:"storagepool"`
	Vm          string `json:"vm" yaml:"vm"`
	Hypervisor  string `json:"hypervisor" yaml:"hypervisor"`
}

//// table: clusters
//type dbClusters struct {
//	CID   uint8  `json:"id" yaml:"id"`
//	Cname string `json:"name" yaml:"name"`
//}
