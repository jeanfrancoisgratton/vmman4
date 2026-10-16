// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/volume_mgt/volumeClone.go

package volume_mgt

import (
	"fmt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// CloneVolume makes a full, independent qcow2 copy of srcVolName into a new
// volume named destVolName, inside the same pool. Unlike a backing-file
// (linked) clone, the result has no runtime dependency on srcVolName -- it
// can be deleted, modified, or rebuilt immediately after cloning with no
// effect on the copy. The caller must Free() the returned volume.
func CloneVolume(conn *libvirt.Connect, poolName, srcVolName, destVolName string) (*libvirt.StorageVol, *ce.CustomError) {
	pool, err := conn.LookupStoragePoolByName(poolName)
	if err != nil {
		return nil, &ce.CustomError{Title: "CloneVolume: pool " + poolName + " not found", Message: err.Error()}
	}
	defer pool.Free()

	srcVol, err := pool.LookupStorageVolByName(srcVolName)
	if err != nil {
		return nil, &ce.CustomError{Title: "CloneVolume: source volume " + srcVolName + " not found in pool " + poolName,
			Message: err.Error()}
	}
	defer srcVol.Free()

	srcInfo, err := srcVol.GetInfo()
	if err != nil {
		return nil, &ce.CustomError{Title: "CloneVolume: cannot get info for " + srcVolName, Message: err.Error()}
	}

	destXML := fmt.Sprintf(`<volume>
  <name>%s</name>
  <capacity unit="bytes">%d</capacity>
  <target>
    <format type="qcow2"/>
  </target>
</volume>`, destVolName, srcInfo.Capacity)

	fmt.Println(hftx.InProgressSign("Cloning "+hftx.Gold(srcVolName)+" to "+hftx.Gold(destVolName)) + ". This might take a while...")

	destVol, err := pool.StorageVolCreateXMLFrom(destXML, srcVol, 0)
	if err != nil {
		return nil, &ce.CustomError{Title: "CloneVolume: cannot clone " + srcVolName + " to " + destVolName, Message: err.Error()}
	}

	if !shared.QuietOutput {
		fmt.Println(hftx.EnabledSign("Volume " + hftx.Gold(srcVolName) + hftx.Green(" cloned " + " as " +
			hftx.Gold(destVolName))))
	}
	return destVol, nil
}
