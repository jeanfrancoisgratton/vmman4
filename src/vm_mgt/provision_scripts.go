// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/provision_scripts.go

package vm_mgt

import (
	"fmt"
	"net"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// buildProvisionScript assembles the guest-side shell script that sets the
// hostname, enforces a single static IPv4 configuration on iface (per env's
// CIDR/gateway), and regenerates SSH host keys. The network-config mechanism
// is selected from osID (as reported live by the guest agent -- see
// guestOSID).
//
// The generated script is meant to be idempotent: every step overwrites or
// removes, none appends, so running it twice leaves the guest in the same
// state. "Enforcing" means that, besides writing our own configuration, it
// removes or neutralises whatever else could configure the same interface
// (other NetworkManager profiles, other netplan files, cloud-init network
// rendering).
func buildProvisionScript(osID, iface, hostname, ipaddr string, env *EnvSpec) (string, *ce.CustomError) {
	ip := net.ParseIP(ipaddr)
	if ip == nil || ip.To4() == nil {
		return "", &ce.CustomError{
			Title:   "buildProvisionScript: only IPv4 addresses are supported",
			Message: ipaddr,
		}
	}
	_, ipnet, err := net.ParseCIDR(env.CIDR)
	if err != nil {
		return "", &ce.CustomError{Title: "buildProvisionScript: invalid 'cidr' in environment file", Message: err.Error()}
	}
	prefix, _ := ipnet.Mask.Size()
	mask := net.IP(ipnet.Mask).String()

	dns4 := ipv4Only(env.DNS)

	// netConfig is the distro-specific part. writeResolv says whether the
	// script itself must write /etc/resolv.conf: it must not where another
	// component owns that file (NetworkManager, netplan/systemd-resolved),
	// since it would be rewritten (or written through a symlink) anyway.
	var netConfig string
	writeResolv := false

	switch osID {
	case "alpine":
		writeResolv = true
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
		writeResolv = true
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
		// netplan: ours is the only YAML left in /etc/netplan. 'gateway4' is
		// deprecated, so the default route is declared explicitly. DNS goes
		// through netplan's 'nameservers' (systemd-resolved owns resolv.conf).
		// 'netplan generate' validates the result without touching the live
		// network; the reboot at the end of ProvisionVM applies it.
		netConfig = fmt.Sprintf(`mkdir -p /etc/netplan
rm -f /etc/netplan/*.yaml
cat > /etc/netplan/01-vmman4.yaml <<'IFACEEOF'
network:
  version: 2
  renderer: networkd
  ethernets:
    %s:
      dhcp4: false
      addresses: [%s/%d]
      routes:
        - to: default
          via: %s
%sIFACEEOF
chown 0:0 /etc/netplan/01-vmman4.yaml
chmod 600 /etc/netplan/01-vmman4.yaml
netplan generate`, iface, ipaddr, prefix, env.Gateway, netplanNameservers(dns4, env.DNSSearch))

	case "rhel", "centos", "fedora", "rocky", "almalinux", "opensuse-tumbleweed", "arch":
		// NetworkManager keyfile. ifcfg files are gone from Fedora and from
		// RHEL 10 and its clones, so this is the only mechanism that works
		// across the whole family. Writing the keyfile directly (rather than
		// calling nmcli) means no dependency on NetworkManager being up when
		// the guest agent answers, and makes the step trivially idempotent:
		// the profile is named after the interface, every other profile is
		// deleted, and NM's auto-generated "Wired connection N" profiles are
		// switched off.
		netConfig = fmt.Sprintf(`NMDIR=/etc/NetworkManager/system-connections
NMFILE="$NMDIR/vmman4-%s.nmconnection"
mkdir -p "$NMDIR" /etc/NetworkManager/conf.d
UUID=$(sed -n 's/^uuid=//p' "$NMFILE" 2>/dev/null | head -n 1)
[ -n "$UUID" ] || UUID=$(cat /proc/sys/kernel/random/uuid)
for f in "$NMDIR"/*; do
  [ -e "$f" ] || continue
  [ "$f" = "$NMFILE" ] || rm -f "$f"
done
cat > /etc/NetworkManager/conf.d/90-vmman4.conf <<'NMCONFEOF'
[main]
no-auto-default=*
NMCONFEOF
cat > "$NMFILE" <<NMEOF
[connection]
id=vmman4-%s
uuid=$UUID
type=ethernet
interface-name=%s
autoconnect=true

[ethernet]

[ipv4]
method=manual
address1=%s/%d,%s
%s
[ipv6]
method=auto
NMEOF
chown 0:0 "$NMFILE" /etc/NetworkManager/conf.d/90-vmman4.conf
chmod 600 "$NMFILE"
chmod 644 /etc/NetworkManager/conf.d/90-vmman4.conf
if command -v restorecon >/dev/null 2>&1; then
  restorecon -R "$NMDIR" /etc/NetworkManager/conf.d || true
fi`, iface, iface, iface, ipaddr, prefix, env.Gateway, nmDNSLines(dns4, env.DNSSearch))

	default:
		return "", &ce.CustomError{
			Title:   "buildProvisionScript: unsupported guest OS",
			Message: fmt.Sprintf("no network provisioning support for OS id %q (interface %s)", osID, iface),
		}
	}

	var resolvBlock string
	if writeResolv {
		var rc strings.Builder
		for _, ns := range env.DNS {
			rc.WriteString("nameserver " + ns + "\n")
		}
		if len(env.DNSSearch) > 0 {
			rc.WriteString("search " + strings.Join(env.DNSSearch, " ") + "\n")
		}
		resolvBlock = fmt.Sprintf("cat > /etc/resolv.conf <<'RESOLVEOF'\n%sRESOLVEOF\n", rc.String())
	}

	return fmt.Sprintf(`#!/bin/sh
set -e

hostnamectl set-hostname %s 2>/dev/null || true
echo %s > /etc/hostname

# cloud-init, if present in the template, would regenerate the network
# configuration (and reset the hostname) at next boot and undo all of this.
if [ -d /etc/cloud/cloud.cfg.d ]; then
  cat > /etc/cloud/cloud.cfg.d/99-vmman4.cfg <<'CLOUDEOF'
network: {config: disabled}
preserve_hostname: true
CLOUDEOF
fi

%s

%s
rm -f /etc/ssh/ssh_host_*_key /etc/ssh/ssh_host_*_key.pub
ssh-keygen -A >/dev/null 2>&1 || true

exit 0
`, hostname, hostname, netConfig, resolvBlock), nil
}

// ipv4Only returns the IPv4 literals of list, in order. It feeds the
// NetworkManager [ipv4] section and netplan's nameserver list; the
// resolv.conf-based branches (alpine, debian) still get the full env.DNS list.
func ipv4Only(list []string) []string {
	var out []string
	for _, s := range list {
		if p := net.ParseIP(s); p != nil && p.To4() != nil {
			out = append(out, s)
		}
	}
	return out
}

// nmDNSLines renders the dns=/dns-search= keys of a NetworkManager keyfile
// [ipv4] section, each terminated by a newline; an empty list yields no line.
// Keyfile lists are ';'-separated with a trailing ';'.
func nmDNSLines(dns, search []string) string {
	var b strings.Builder
	if len(dns) > 0 {
		b.WriteString("dns=" + strings.Join(dns, ";") + ";\n")
		b.WriteString("ignore-auto-dns=true\n")
	}
	if len(search) > 0 {
		b.WriteString("dns-search=" + strings.Join(search, ";") + ";\n")
	}
	return b.String()
}

// netplanNameservers renders the 'nameservers:' block of a netplan ethernet
// stanza (indented to sit under the interface name), newline-terminated; an
// empty dns list yields an empty string.
func netplanNameservers(dns, search []string) string {
	if len(dns) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("      nameservers:\n")
	b.WriteString("        addresses: [" + strings.Join(dns, ", ") + "]\n")
	if len(search) > 0 {
		b.WriteString("        search: [" + strings.Join(search, ", ") + "]\n")
	}
	return b.String()
}
