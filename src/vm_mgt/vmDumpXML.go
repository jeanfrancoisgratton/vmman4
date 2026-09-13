// vmman3 : Écrit par Jean-François Gratton (jean-francois@famillegratton.net)
// src/vmmanagement/vmDumpXML.go
// 2022-11-05 13:45:39

package vm_mgt

import (
	"os"
	"strings"
	"vmman4/connection_mgt"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// XmlDump() : dumps the vm config in an xml file
// libvirt call: https://pkg.go.dev/libvirt.org/go/libvirt#Domain.GetXMLDesc
// TODO: more robust error handling here...

func DumpVmXML(vmname string, xmlfile string) *ce.CustomError {
	var conn *libvirt.Connect
	var err *ce.CustomError

	if !strings.HasSuffix(xmlfile, ".xml") {
		xmlfile += ".xml"
	}

	if err = connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	if conn, err = shared.Connect2HVM(); err != nil {
		return err
	}
	defer conn.Close()

	domain, err := shared.GetDomain(conn, vmname)
	if err != nil {
		return err
	}
	defer domain.Free()

	// Shut the VM down, if active
	shared.Wait4Shutdown(domain, vmname)
	data, _ := domain.GetXMLDesc(libvirt.DOMAIN_XML_SECURE | libvirt.DOMAIN_XML_INACTIVE | libvirt.DOMAIN_XML_MIGRATABLE)

	file, e := os.Create(xmlfile)
	if e != nil {
		return &ce.CustomError{Title: "os.Create() failed: ", Message: e.Error()}
	}
	defer file.Close()
	if _, e := file.WriteString(data); e != nil {
		return &ce.CustomError{Title: "file.WriteString() : Unable to write XML file", Message: e.Error()}
	}
	file.Sync()

	return nil
}
