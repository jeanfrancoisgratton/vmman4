// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vmmanagement/types.go
// Original timestamp: 2026/05/27 18:46:59

package vmmanagement

var ForceConsoleConnection = false

type vmInfo struct {
	viId              uint
	viName            string
	viState           string
	viMem             uint64
	viCpu             uint
	viSnapshots       uint
	viCurrentSnapshot string
	viInterfaceName   string
	viIPaddress       string
	//viHypervisor       string
	//viOperatingSystem  string
	//viLastStatusChange string
	//viStoragePool      string
}
