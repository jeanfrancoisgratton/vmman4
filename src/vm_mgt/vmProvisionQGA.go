// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/vmProvisionQGA.go

package vm_mgt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

type qgaRequest struct {
	Execute   string      `json:"execute"`
	Arguments interface{} `json:"arguments,omitempty"`
}

// qgaCall sends a single QGA command and, if out is non-nil, unmarshals the
// response into it.
func qgaCall(dom *libvirt.Domain, execute string, args interface{}, out interface{}) *ce.CustomError {
	reqJSON, err := json.Marshal(qgaRequest{Execute: execute, Arguments: args})
	if err != nil {
		return &ce.CustomError{Title: "qgaCall: cannot marshal request", Message: err.Error()}
	}

	respJSON, err := dom.QemuAgentCommand(string(reqJSON), libvirt.DOMAIN_QEMU_AGENT_COMMAND_BLOCK, 0)
	if err != nil {
		return &ce.CustomError{Title: fmt.Sprintf("qgaCall: %s failed", execute), Message: err.Error()}
	}

	if out != nil {
		if err := json.Unmarshal([]byte(respJSON), out); err != nil {
			return &ce.CustomError{Title: fmt.Sprintf("qgaCall: cannot parse %s response", execute), Message: err.Error()}
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

// ─── guest-file-* / guest-exec ─────────────────────────────────────────────────

type qgaFileOpenArgs struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
}
type qgaFileOpenResult struct {
	Return int `json:"return"`
}

type qgaFileWriteArgs struct {
	Handle int    `json:"handle"`
	BufB64 string `json:"buf-b64"`
}

type qgaFileCloseArgs struct {
	Handle int `json:"handle"`
}

// qgaWriteFile writes content to path inside the guest, via QGA (no network,
// no SSH -- just the virtio-serial channel).
func qgaWriteFile(dom *libvirt.Domain, path, content string) *ce.CustomError {
	var openRes qgaFileOpenResult
	if cerr := qgaCall(dom, "guest-file-open", qgaFileOpenArgs{Path: path, Mode: "w+"}, &openRes); cerr != nil {
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

type qgaExecArgs struct {
	Path          string   `json:"path"`
	Arg           []string `json:"arg,omitempty"`
	CaptureOutput bool     `json:"capture-output"`
}
type qgaExecResult struct {
	Return struct {
		PID int `json:"pid"`
	} `json:"return"`
}

type qgaExecStatusArgs struct {
	PID int `json:"pid"`
}
type qgaExecStatusResult struct {
	Return struct {
		Exited   bool   `json:"exited"`
		Exitcode int    `json:"exitcode"`
		OutData  string `json:"out-data"`
		ErrData  string `json:"err-data"`
	} `json:"return"`
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
