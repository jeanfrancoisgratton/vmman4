package main

import (
	"os"
	"path/filepath"

	"vmman4/cmd"
)

func main() {
	//fmt.Printf("\x1bc")
	_ = os.Mkdir(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4"), os.ModePerm)

	cmd.Execute()
}
