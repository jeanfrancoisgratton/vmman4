// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmProvisionScripts.go

package vm_mgt

import (
	"fmt"
	"net"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// buildProvisionScript assembles the guest-side shell script that sets the
// hostname, configures iface with a static ipaddr (per env's CIDR/gateway),
// writes resolv.conf, and regenerates SSH host keys. The network-config
// mechanism is selected from osID (as reported live by the guest agent --
// see guestOSID), mirroring vmman1's per-distro handling.
func buildProvisionScript(osID, iface, hostname, ipaddr string, env *EnvSpec) (string, *ce.CustomError) {
	_, ipnet, err := net.ParseCIDR(env.CIDR)
	if err != nil {
		return "", &ce.CustomError{Title: "buildProvisionScript: invalid 'cidr' in environment file", Message: err.Error()}
	}
	prefix, _ := ipnet.Mask.Size()
	mask := net.IP(ipnet.Mask).String()

	var resolvConf strings.Builder
	for _, ns := range env.DNS {
		resolvConf.WriteString("nameserver " + ns + "\n")
	}
	if len(env.DNSSearch) > 0 {
		resolvConf.WriteString("search " + strings.Join(env.DNSSearch, " ") + "\n")
	}

	var netConfig string
	switch osID {
	case "alpine":
		netConfig = fmt.Sprintf(`rm -f /etc/network/interfaces
cat > /etc/network/interfaces <<'IFACEEOF'
auto lo
iface lo inet loopback

auto %s
iface %s inet static
	address %s
	netmask %s
	gateway %s
IFACEEOF`, iface, iface, ipaddr, mask, env.Gateway)

	case "debian":
		netConfig = fmt.Sprintf(`cat > /etc/network/interfaces <<'IFACEEOF'
source /etc/network/interfaces.d/*

auto lo
iface lo inet loopback

allow-hotplug %s
iface %s inet static
	address %s/%d
	gateway %s
IFACEEOF
chown 0:0 /etc/network/interfaces
chmod 644 /etc/network/interfaces`, iface, iface, ipaddr, prefix, env.Gateway)

	case "ubuntu":
		netConfig = fmt.Sprintf(`rm -f /etc/netplan/*.yaml
cat > /etc/netplan/01-vmman4.yaml <<'IFACEEOF'
network:
  version: 2
  ethernets:
    %s:
      addresses: [%s/%d]
      gateway4: %s
IFACEEOF
netplan apply`, iface, ipaddr, prefix, env.Gateway)

	case "rhel", "centos", "fedora", "rocky", "almalinux":
		netConfig = fmt.Sprintf(`egrep -v "IPADDR|GATEWAY|PREFIX|BOOTPROTO" /etc/sysconfig/network-scripts/ifcfg-%s > /tmp/.vmman4-ifcfg 2>/dev/null || true
{
  echo BOOTPROTO=static
  echo IPADDR=%s
  echo GATEWAY=%s
  echo PREFIX=%d
  cat /tmp/.vmman4-ifcfg 2>/dev/null
} > /etc/sysconfig/network-scripts/ifcfg-%s
rm -f /tmp/.vmman4-ifcfg
chown 0:0 /etc/sysconfig/network-scripts/ifcfg-%s
chmod 644 /etc/sysconfig/network-scripts/ifcfg-%s`, iface, ipaddr, env.Gateway, prefix, iface, iface, iface)

	case "opensuse-tumbleweed", "arch":
		// Both templates are NetworkManager-managed: resolve iface's actual
		// connection profile (rather than assuming a name like "Wired connection 1")
		// and reconfigure that profile in place.
		netConfig = fmt.Sprintf(`CONN=$(nmcli -g GENERAL.CONNECTION device show %s)
nmcli con mod "$CONN" ipv4.addresses %s/%d ipv4.gateway %s ipv4.dns "%s" ipv4.dns-search "%s" ipv4.method manual
nmcli con up "$CONN"`, iface, ipaddr, prefix, env.Gateway, strings.Join(env.DNS, ","), strings.Join(env.DNSSearch, ","))

	default:
		return "", &ce.CustomError{
			Title:   "buildProvisionScript: unsupported guest OS",
			Message: fmt.Sprintf("no network provisioning support for OS id %q (interface %s)", osID, iface),
		}
	}

	return fmt.Sprintf(`#!/bin/sh
set -e

hostnamectl set-hostname %s 2>/dev/null || true
echo %s > /etc/hostname

%s

cat > /etc/resolv.conf <<'RESOLVEOF'
%sRESOLVEOF

rm -f /etc/ssh/ssh_host_*_key /etc/ssh/ssh_host_*_key.pub
ssh-keygen -A >/dev/null 2>&1 || true

exit 0
`, hostname, hostname, netConfig, resolvConf.String()), nil
}
