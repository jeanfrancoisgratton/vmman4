// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/provision_qemu-guest-agent.go

package vm_mgt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"libvirt.org/go/libvirt"
)

// ─── Raw QEMU guest agent plumbing ─────────────────────────────────────────────
//
// The Go binding only exposes Domain.QemuAgentCommand, a pass-through for raw
// QGA JSON-RPC (there are no typed helpers for guest-exec/guest-file-* the way
// there are for, say, domain lifecycle calls), so every command below is
// hand-built JSON.

// qgaCall sends a single QGA command and, if out is non-nil, unmarshals the
// response into it.
func qgaCall(dom *libvirt.Domain, execute string, args interface{}, out interface{}) *ce.CustomError {
	reqJSON, err := json.Marshal(qgaRequest{Execute: execute, Arguments: args})
	if err != nil {
		return &ce.CustomError{Title: "qgaCall: cannot marshal request", Message: err.Error()}
	}

	respJSON, err := dom.QemuAgentCommand(string(reqJSON), libvirt.DOMAIN_QEMU_AGENT_COMMAND_BLOCK, 0)
	if err != nil {
		return &ce.CustomError{Title: "qgaCall: " + execute + " failed", Message: err.Error()}
	}

	if out != nil {
		if err := json.Unmarshal([]byte(respJSON), out); err != nil {
			return &ce.CustomError{Title: "qgaCall: cannot parse " + execute + " response", Message: err.Error()}
		}
	}
	return nil
}

// waitForGuestAgent polls guest-ping until the agent answers or timeout elapses.
func waitForGuestAgent(dom *libvirt.Domain, timeout time.Duration) *ce.CustomError {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		_, err := dom.QemuAgentCommand(`{"execute":"guest-ping"}`, libvirt.DOMAIN_QEMU_AGENT_COMMAND_DEFAULT, 0)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(3 * time.Second)
	}
	msg := "timed out"
	if lastErr != nil {
		msg = lastErr.Error()
	}
	return &ce.CustomError{Title: "waitForGuestAgent: QEMU guest agent never became reachable", Message: msg}
}

// guestOSID returns the guest's /etc/os-release ID field (e.g. "ubuntu",
// "debian", "rhel", "alpine"), as reported live by the guest agent.
func guestOSID(dom *libvirt.Domain) (string, *ce.CustomError) {
	info, err := dom.GetGuestInfo(libvirt.DOMAIN_GUEST_INFO_OS, 0)
	if err != nil {
		return "", &ce.CustomError{Title: "guestOSID: GetGuestInfo failed", Message: err.Error()}
	}
	if info.OS == nil || !info.OS.IDSet || info.OS.ID == "" {
		return "", &ce.CustomError{Title: "guestOSID: guest agent did not report an OS id"}
	}
	return info.OS.ID, nil
}

// guestPrimaryInterface returns the first non-loopback NIC name the guest
// agent reports (e.g. "ens3", "eth0") -- the real, kernel-assigned name, not
// a guess.
func guestPrimaryInterface(dom *libvirt.Domain) (string, *ce.CustomError) {
	info, err := dom.GetGuestInfo(libvirt.DOMAIN_GUEST_INFO_INTERFACES, 0)
	if err != nil {
		return "", &ce.CustomError{Title: "guestPrimaryInterface: GetGuestInfo failed", Message: err.Error()}
	}
	for _, iface := range info.Interfaces {
		if !iface.NameSet || iface.Name == "lo" {
			continue
		}
		return iface.Name, nil
	}
	return "", &ce.CustomError{Title: "guestPrimaryInterface: guest agent reported no non-loopback interface"}
}

// ─── Guest agent capability preflight ──────────────────────────────────────────
//
// Several distros ship qemu-ga with the guest-file-* and guest-exec* commands
// blocked by default (the RHEL family via /etc/sysconfig/qemu-ga, openSUSE
// similarly). Calling a blocked command fails with a terse "has been disabled"
// error, possibly after other work was done, so guest-info -- which lists every
// command along with its enabled flag and is never blocked -- is consulted
// first.

// qgaRequiredCommands are the agent commands runProvisionScript relies on.
var qgaRequiredCommands = []string{
	"guest-file-open",
	"guest-file-write",
	"guest-file-close",
	"guest-exec",
	"guest-exec-status",
}

type qgaInfoResult struct {
	Return struct {
		Version           string `json:"version"`
		SupportedCommands []struct {
			Name    string `json:"name"`
			Enabled bool   `json:"enabled"`
		} `json:"supported_commands"`
	} `json:"return"`
}

// missingQGACommands returns those of required that info reports as disabled
// or does not list at all (an agent too old to know the command), in the order
// of required.
func missingQGACommands(info *qgaInfoResult, required []string) []string {
	enabled := make(map[string]bool, len(info.Return.SupportedCommands))
	for _, c := range info.Return.SupportedCommands {
		enabled[c.Name] = c.Enabled
	}
	var missing []string
	for _, name := range required {
		if !enabled[name] {
			missing = append(missing, name)
		}
	}
	return missing
}

// checkGuestAgentCommands fails early, and explicitly, when the template's
// guest agent blocks any command the provisioning needs.
func checkGuestAgentCommands(dom *libvirt.Domain) *ce.CustomError {
	var info qgaInfoResult
	if cerr := qgaCall(dom, "guest-info", nil, &info); cerr != nil {
		return cerr
	}
	missing := missingQGACommands(&info, qgaRequiredCommands)
	if len(missing) == 0 {
		return nil
	}
	return &ce.CustomError{
		Title: "checkGuestAgentCommands: the template's QEMU guest agent blocks commands that vm provision needs",
		Message: fmt.Sprintf("disabled or unsupported: %s (qemu-ga %s). This is fixed in the template, not here: "+
			"find where its block/allow list is set (see 'systemctl cat qemu-guest-agent' and /etc/sysconfig/qemu-ga; "+
			"typically a BLACKLIST_RPC, BLOCK_RPCS or FILTER_RPC_ARGS line), remove or comment it out, restart "+
			"qemu-guest-agent, and re-clone from the updated template",
			strings.Join(missing, ", "), info.Return.Version),
	}
}

// ─── guest-file-* / guest-exec ─────────────────────────────────────────────────

// qgaWriteFile writes content to path inside the guest, via QGA (no network,
// no SSH -- just the virtio-serial channel).
func qgaWriteFile(dom *libvirt.Domain, path, content string) *ce.CustomError {
	var openRes qgaFileOpenResult
	if cerr := qgaCall(dom, "guest-file-open", qgaFileOpenArgs{Path: path, Mode: "xw+"}, &openRes); cerr != nil {
		return cerr
	}
	handle := openRes.Return

	writeArgs := qgaFileWriteArgs{Handle: handle, BufB64: base64.StdEncoding.EncodeToString([]byte(content))}
	writeErr := qgaCall(dom, "guest-file-write", writeArgs, nil)

	// Always attempt to close the handle, even if the write failed.
	closeErr := qgaCall(dom, "guest-file-close", qgaFileCloseArgs{Handle: handle}, nil)

	if writeErr != nil {
		return writeErr
	}
	return closeErr
}

// qgaExec runs path (with args) inside the guest via guest-exec, polling
// guest-exec-status until completion, and returns its exit code and combined
// stdout+stderr.
func qgaExec(dom *libvirt.Domain, path string, args []string, timeout time.Duration) (int, string, *ce.CustomError) {
	var execRes qgaExecResult
	execArgs := qgaExecArgs{Path: path, Arg: args, CaptureOutput: true}
	if cerr := qgaCall(dom, "guest-exec", execArgs, &execRes); cerr != nil {
		return 0, "", cerr
	}
	pid := execRes.Return.PID

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		var statusRes qgaExecStatusResult
		if cerr := qgaCall(dom, "guest-exec-status", qgaExecStatusArgs{PID: pid}, &statusRes); cerr != nil {
			return 0, "", cerr
		}
		if statusRes.Return.Exited {
			outBytes, _ := base64.StdEncoding.DecodeString(statusRes.Return.OutData)
			errBytes, _ := base64.StdEncoding.DecodeString(statusRes.Return.ErrData)
			return statusRes.Return.Exitcode, string(outBytes) + string(errBytes), nil
		}
		time.Sleep(1 * time.Second)
	}
	return 0, "", &ce.CustomError{
		Title:   "qgaExec: guest command timed out",
		Message: fmt.Sprintf("path=%s pid=%d", path, pid),
	}
}

// runProvisionScript pushes script into the guest and executes it via /bin/sh.
func runProvisionScript(dom *libvirt.Domain, script string) *ce.CustomError {
	const scriptPath = "/tmp/.vmman4-provision.sh"

	if cerr := checkGuestAgentCommands(dom); cerr != nil {
		return cerr
	}

	if cerr := qgaWriteFile(dom, scriptPath, script); cerr != nil {
		return cerr
	}

	exitcode, out, cerr := qgaExec(dom, "/bin/sh", []string{scriptPath}, 60*time.Second)
	if cerr != nil {
		return cerr
	}
	if exitcode != 0 {
		return &ce.CustomError{
			Title:   "runProvisionScript: guest-side provisioning script failed",
			Message: fmt.Sprintf("exit=%d output=%s", exitcode, out),
		}
	}

	// Best-effort cleanup; a leftover /tmp file is not worth failing the run over.
	_, _, _ = qgaExec(dom, "/bin/rm", []string{"-f", scriptPath}, 10*time.Second)
	return nil
}
