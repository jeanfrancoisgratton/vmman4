// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/vm_mgt/types_provision.go

package vm_mgt

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	ce "github.com/jeanfrancoisgratton/customError/v3"
)

// EnvFile is the -E/--environment flag; empty means "use DefaultEnvFile()".
var EnvFile = ""

type qgaRequest struct {
	Execute   string `json:"execute"`
	Arguments any    `json:"arguments,omitempty"`
}

// EnvSpec is unmarshalled from the environment JSON file 'vm provision' reads
// network defaults from.
type EnvSpec struct {
	PoolName  string   `json:"poolname"`  // libvirt storage pool holding VM disks
	CIDR      string   `json:"cidr"`      // network block, e.g. "10.0.0.0/8"
	Gateway   string   `json:"gateway"`   // default gateway
	DNS       []string `json:"dns"`       // one or more nameservers
	DNSSearch []string `json:"dnssearch"` // DNS search domains, may be empty
}

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

// DefaultEnvFile returns ~/.config/JFG/vmman4/env.json, the default
// environment file path used when -E is not passed.
func DefaultEnvFile() string {
	return filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", "env.json")
}

// loadEnv reads and validates the environment file.
func loadEnv(envFile string) (*EnvSpec, *ce.CustomError) {
	data, err := os.ReadFile(envFile)
	if err != nil {
		return nil, &ce.CustomError{Title: "loadEnv: cannot read file", Message: err.Error()}
	}

	var env EnvSpec
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, &ce.CustomError{Title: "loadEnv: invalid JSON", Message: err.Error()}
	}

	if env.PoolName == "" {
		return nil, &ce.CustomError{Title: "loadEnv: 'poolname' is required"}
	}
	if _, _, err := net.ParseCIDR(env.CIDR); err != nil {
		return nil, &ce.CustomError{Title: "loadEnv: invalid 'cidr'", Message: err.Error()}
	}
	if net.ParseIP(env.Gateway) == nil {
		return nil, &ce.CustomError{Title: "loadEnv: invalid 'gateway'", Message: env.Gateway}
	}
	if len(env.DNS) == 0 {
		return nil, &ce.CustomError{Title: "loadEnv: 'dns' must list at least one nameserver"}
	}

	return &env, nil
}

// checkIPInCIDR ensures ipaddr is a valid address that falls within env.CIDR.
func checkIPInCIDR(ipaddr string, env *EnvSpec) *ce.CustomError {
	ip := net.ParseIP(ipaddr)
	if ip == nil {
		return &ce.CustomError{Title: "checkIPInCIDR: invalid IP address", Message: ipaddr}
	}
	_, ipnet, err := net.ParseCIDR(env.CIDR)
	if err != nil {
		return &ce.CustomError{Title: "checkIPInCIDR: invalid 'cidr' in environment file", Message: err.Error()}
	}
	if !ipnet.Contains(ip) {
		return &ce.CustomError{
			Title:   "checkIPInCIDR: IP address outside the configured network",
			Message: fmt.Sprintf("%s is not within %s", ipaddr, env.CIDR),
		}
	}
	return nil
}
