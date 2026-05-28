// vmman3 : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/vmmanagement/vmRename.go
// 2022-10-29 19:04:58

package vmmanagement

import (
	"fmt"
	"strconv"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
	"vmman4/connection"
	"vmman4/shared"
)

func Rename(oldName, newName string) *ce.CustomError {
	var snapshotflags libvirt.DomainSnapshotListFlags
	var err *ce.CustomError
	var conn *libvirt.Connect
	var domain *libvirt.Domain

	if err = connection.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = shared.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	if domain, err = shared.GetDomain(conn, oldName); err != nil {
		return err
	}
	defer domain.Free()

	// The domain cannot be renamed if there are any snapshot
	numsnap, _ := domain.SnapshotNum(snapshotflags)
	if numsnap > 0 {
		fmt.Println()
		fmt.Println(hftx.EuropeanStopSign("VM " + hftx.Blue(oldName) + " holds " + hftx.Red(strconv.Itoa(numsnap)) + " snapshot and cannot proceed\n"))
		return &ce.CustomError{Fatality: ce.Continuable, Title: "Cannot rename VM " + oldName + " to " + newName,
			Message: "You cannot rename " + oldName + " as this VM holds snapshots. The snapshots need to be removed, first."}
	}

	Wait4Shutdown(domain, oldName)
	if serr := domain.Rename(newName, 0); serr != nil {
		return &ce.CustomError{Title: "Cannot rename " + oldName + " to " + newName, Message: serr.Error()}
	}
	fmt.Println(hftx.EnabledSign(oldName + hftx.Green(" --> ") + newName))

	return nil
}
