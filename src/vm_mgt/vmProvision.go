// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmProvision.go

package vm_mgt

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"vmman4/connection_mgt"
	"vmman4/shared"
	"vmman4/volume_mgt"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"libvirt.org/go/libvirt"
)

// ProvisionVM clones templateName (an existing, already-defined domain whose
// disk is templateName.qcow2) into a brand-new VM named hostname, assigns it
// ipaddr (validated against the environment file's CIDR), boots it, and --
// via the QEMU guest agent, assumed to already be installed and running in
// the template -- sets the guest's hostname and network configuration and
// regenerates its SSH host keys (cloning the disk also clones the template's
// host keys, so every clone must get fresh ones).
//
// envFile is the environment JSON path; pass "" to use DefaultEnvFile().
func ProvisionVM(hostname, ipaddr, templateName, envFile string) *ce.CustomError {
	if envFile == "" {
		envFile = DefaultEnvFile()
	}
	env, cerr := loadEnv(envFile)
	if cerr != nil {
		return cerr
	}
	if cerr := checkIPInCIDR(ipaddr, env); cerr != nil {
		return cerr
	}

	if err := connection_mgt.ResolveConnectionURI(); err != nil {
		return err
	}
	conn, err := shared.Connect2HVM()
	if err != nil {
		return err
	}
	defer conn.Close()

	destVol, cerr := volume_mgt.CloneVolume(conn, env.PoolName, templateName+".qcow2", hostname+".qcow2")
	if cerr != nil {
		return cerr
	}
	defer destVol.Free()

	destPath, gerr := destVol.GetPath()
	if gerr != nil {
		return &ce.CustomError{Title: "ProvisionVM: cannot get path of cloned volume", Message: gerr.Error()}
	}
	fmt.Println(hftx.EnabledSign(hostname + ".qcow2" + hftx.Green(" cloned from "+templateName+".qcow2")))

	dom, cerr := cloneDomainDefinition(conn, env.PoolName, templateName, hostname, destPath)
	if cerr != nil {
		return cerr
	}
	defer dom.Free()
	fmt.Println(hftx.EnabledSign(hostname + hftx.Green(" DEFINED")))

	if err := dom.Create(); err != nil {
		return &ce.CustomError{Title: "ProvisionVM: cannot start " + hostname, Message: err.Error()}
	}
	fmt.Println(hftx.EnabledSign(hostname + hftx.Green(" started, waiting for QEMU guest agent")))

	if cerr := waitForGuestAgent(dom, 2*time.Minute); cerr != nil {
		return cerr
	}

	osID, cerr := guestOSID(dom)
	if cerr != nil {
		return cerr
	}
	ifaceName, cerr := guestPrimaryInterface(dom)
	if cerr != nil {
		return cerr
	}

	script, cerr := buildProvisionScript(osID, ifaceName, hostname, ipaddr, env)
	if cerr != nil {
		return cerr
	}
	if cerr := runProvisionScript(dom, script); cerr != nil {
		return cerr
	}

	if err := dom.Reboot(0); err != nil {
		return &ce.CustomError{Title: "ProvisionVM: reboot failed", Message: err.Error()}
	}

	fmt.Println(hftx.EnabledSign(hostname + hftx.Green(fmt.Sprintf(" PROVISIONED (%s, %s on %s), rebooting to apply network configuration", osID, ipaddr, ifaceName))))
	return nil
}

var (
	domNameRe = regexp.MustCompile(`(?s)<name>.*?</name>`)
	domUUIDRe = regexp.MustCompile(`(?s)\s*<uuid>.*?</uuid>`)
	domMACRe  = regexp.MustCompile(`(?s)\s*<mac address=["'][^"']*["']\s*/>`)
)

// cloneDomainDefinition looks up templateName's domain and disk, and defines
// a new domain named newName whose disk is newDiskPath: the template's raw
// XML is reused verbatim except for its <name> (renamed), <uuid> and
// <mac address> elements (stripped -- libvirt assigns fresh ones on define),
// and its disk <source> (repointed at the cloned volume).
func cloneDomainDefinition(conn *libvirt.Connect, poolName, templateName, newName, newDiskPath string) (*libvirt.Domain, *ce.CustomError) {
	pool, err := conn.LookupStoragePoolByName(poolName)
	if err != nil {
		return nil, &ce.CustomError{Title: fmt.Sprintf("cloneDomainDefinition: pool %q not found", poolName), Message: err.Error()}
	}
	defer pool.Free()

	srcVol, err := pool.LookupStorageVolByName(templateName + ".qcow2")
	if err != nil {
		return nil, &ce.CustomError{
			Title:   fmt.Sprintf("cloneDomainDefinition: template volume %q not found in pool %q", templateName+".qcow2", poolName),
			Message: err.Error(),
		}
	}
	defer srcVol.Free()
	oldDiskPath, err := srcVol.GetPath()
	if err != nil {
		return nil, &ce.CustomError{Title: "cloneDomainDefinition: cannot get template volume path", Message: err.Error()}
	}

	tmplDom, err := conn.LookupDomainByName(templateName)
	if err != nil {
		return nil, &ce.CustomError{Title: fmt.Sprintf("cloneDomainDefinition: template domain %q not found", templateName), Message: err.Error()}
	}
	defer tmplDom.Free()

	xmlDesc, err := tmplDom.GetXMLDesc(0)
	if err != nil {
		return nil, &ce.CustomError{Title: "cloneDomainDefinition: cannot get template domain XML", Message: err.Error()}
	}

	newXML := domNameRe.ReplaceAllString(xmlDesc, "<name>"+newName+"</name>")
	newXML = domUUIDRe.ReplaceAllString(newXML, "")
	newXML = domMACRe.ReplaceAllString(newXML, "")
	newXML = strings.ReplaceAll(newXML, oldDiskPath, newDiskPath)

	dom, err := conn.DomainDefineXML(newXML)
	if err != nil {
		return nil, &ce.CustomError{Title: "cloneDomainDefinition: DomainDefineXML failed", Message: err.Error()}
	}
	return dom, nil
}
