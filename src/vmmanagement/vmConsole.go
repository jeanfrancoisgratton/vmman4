// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vmmanagement/vmConsole.go
// Original timestamp: 2026/05/27

package vmmanagement

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	"golang.org/x/term"
	"libvirt.org/go/libvirt"
	"vmman4/connection"
	"vmman4/shared"
)

// escapeChar is the session escape sequence (Ctrl+]), matching virsh behaviour
const escapeChar = 0x1d

// Console opens an interactive serial console session on the given VM,
// equivalent to `virsh console [--force]`.
// When ForceConsoleConnection is true (the -f flag), any existing console
// session is forcefully evicted before connecting.
func Console(vmname string) *ce.CustomError {
	if err := connection.ResolveConnectionURI(); err != nil {
		return err
	}

	conn, cerr := shared.Connect2HVM()
	if cerr != nil {
		return cerr
	}
	defer conn.Close()

	domain, cerr := shared.GetDomain(conn, vmname)
	if cerr != nil {
		return cerr
	}
	defer domain.Free()

	active, aerr := domain.IsActive()
	if aerr != nil {
		return &ce.CustomError{Title: "Console error", Message: aerr.Error()}
	}
	if !active {
		return &ce.CustomError{
			Title:   "Console error",
			Message: "domain '" + vmname + "' is not running",
		}
	}

	// A blocking stream (flags=0) is what we want: Recv() will sleep until
	// data arrives, which is the right model for an interactive console.
	stream, serr := conn.NewStream(0)
	if serr != nil {
		return &ce.CustomError{Title: "Failed to create stream", Message: serr.Error()}
	}
	defer stream.Free()

	// Mirror `virsh console`: pass FORCE when -f is set, 0 otherwise.
	// DOMAIN_CONSOLE_SAFE is intentionally not used as the default — it
	// prevents connecting to an unattended console even when no peer is
	// present, which is overly restrictive for a management tool.
	var consoleFlags libvirt.DomainConsoleFlags
	if ForceConsoleConnection {
		consoleFlags = libvirt.DOMAIN_CONSOLE_FORCE
	}

	// Empty devname → libvirt picks the first/default serial console.
	if oerr := domain.OpenConsole("", stream, consoleFlags); oerr != nil {
		return &ce.CustomError{
			Title:   "Failed to open console on '" + vmname + "'",
			Message: oerr.Error(),
		}
	}

	// Raw terminal mode is required so keystrokes are forwarded immediately
	// and control characters (arrows, Ctrl-C …) pass through to the guest.
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return &ce.CustomError{
			Title:   "Console error",
			Message: "stdin is not a terminal — cannot open interactive console",
		}
	}
	oldState, terr := term.MakeRaw(fd)
	if terr != nil {
		return &ce.CustomError{Title: "Failed to set raw terminal mode", Message: terr.Error()}
	}
	defer term.Restore(fd, oldState) //nolint:errcheck

	// \r\n required in raw mode: \n alone won't return to column 0.
	fmt.Fprintf(os.Stderr, "Connected to domain '%s'\r\nEscape character is ^] (Ctrl+])\r\n", vmname)

	// quit carries the first termination cause from either goroutine.
	quit := make(chan error, 2)

	// --- guest → local stdout ---
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stream.Recv(buf)
			if n > 0 {
				os.Stdout.Write(buf[:n]) //nolint:errcheck
			}
			if err != nil {
				// io.EOF means the guest closed the stream normally.
				quit <- err
				return
			}
		}
	}()

	// --- local stdin → guest (with Ctrl+] escape detection) ---
	go func() {
		buf := make([]byte, 256)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				quit <- err
				return
			}
			// Scan for the escape byte before forwarding.
			for i := 0; i < n; i++ {
				if buf[i] == escapeChar {
					quit <- nil
					return
				}
			}
			if n > 0 {
				if _, werr := stream.Send(buf[:n]); werr != nil {
					quit <- werr
					return
				}
			}
		}
	}()

	// Catch SIGTERM / SIGINT so Ctrl-C on the host (not inside the guest)
	// exits cleanly instead of leaving the terminal in raw mode.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigCh)

	select {
	case err := <-quit:
		switch {
		case err == nil:
			fmt.Fprintf(os.Stderr, "\r\nEscape sequence received, disconnecting from '%s'\r\n", vmname)
		case err == io.EOF:
			fmt.Fprintf(os.Stderr, "\r\nConsole stream closed by guest\r\n")
		default:
			fmt.Fprintf(os.Stderr, "\r\nConsole stream error: %v\r\n", err)
		}
		stream.Abort() //nolint:errcheck

	case sig := <-sigCh:
		fmt.Fprintf(os.Stderr, "\r\nSignal %v received, disconnecting from '%s'\r\n", sig, vmname)
		stream.Abort() //nolint:errcheck
	}

	return nil
}
