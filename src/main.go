package main

import (
	"os"
	"path/filepath"

	"vmman4/cmd"
)

// envSample is a ready-to-use baseline for ~/.config/JFG/vmman4/env.json, the
// environment file 'vm provision' reads network defaults (pool, CIDR, gateway,
// DNS) from. Refreshed on every run so it always reflects the current schema;
// env.json itself is never touched here.
const envSample = `{
  "poolname": "vmpool",
  "cidr": "10.0.0.0/8",
  "gateway": "10.1.1.1",
  "dns": ["8.8.8.8", "1.1.1.1"],
  "dnssearch": ["example.com"]
}
`

func main() {
	//fmt.Printf("\x1bc")
	confDir := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4")
	_ = os.Mkdir(confDir, os.ModePerm)
	_ = os.WriteFile(filepath.Join(confDir, "env-sample.json"), []byte(envSample), 0644)

	cmd.Execute()
}
