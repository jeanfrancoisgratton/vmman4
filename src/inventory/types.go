// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/inventory/types.go
// Original timestamp: 2026/05/22 07:55:18

package inventory

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
