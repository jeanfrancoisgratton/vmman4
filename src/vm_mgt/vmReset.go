// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmReset.go
// Original timestamp : 2026.09.11 16:43:26

package vm_mgt

import ce "github.com/jeanfrancoisgratton/customError/v3"

// This one is a simple one, it wraps ResetAll() over StopAll() and StartAll()

func ResetAllVMs() *ce.CustomError {
	if e := StopAll(); e != nil {
		return e
	}

	return StartAll()
}

// Same principle here: we wrap Reset() around Stop() and Start()

func ResetVM(vmlist []string) *ce.CustomError {

	if e := StopVM(vmlist); e != nil {
		return e
	}

	return StartVM(vmlist)
}
