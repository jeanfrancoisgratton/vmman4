// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/types.go
// Original timestamp: 2026/05/27 18:46:59

package vm_mgt

var ForceConsoleConnection = false
var KeepStorage = false
var SampleMode = false

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
	viDiskCount       uint
	viDiskTotalSize   uint64 // bytes
	//viHypervisor       string
	//viOperatingSystem  string
	//viLastStatusChange string
	//viStoragePool string
	//viDisks       []volume_mgt.DiskInfo
}
