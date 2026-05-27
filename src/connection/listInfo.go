// vmman4
// Written by J.F.Gratton <jean-francois@famillegratton.net>
// Original filename: src/connection/listInfo.go
// Original timestamp: 2026/05/20 08:09:45

package connection

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func ConnList() *ce.CustomError {
	var err error
	var dirFH *os.File
	var finfo, fileInfos []os.FileInfo

	// list environment files
	envdir := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4")

	if dirFH, err = os.Open(envdir); err != nil {
		return &ce.CustomError{Title: "Unable to read config directory", Message: err.Error()}
	}
	defer dirFH.Close()

	if fileInfos, err = dirFH.Readdir(0); err != nil {
		return &ce.CustomError{Title: "Unable to read files in config directory", Message: err.Error()}
	}

	for _, info := range fileInfos {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			finfo = append(finfo, info)
		}
	}

	//if err != nil {
	//	return err
	//}

	fmt.Printf("Number of environment files: %s\n", hftx.Green(fmt.Sprintf("%d", len(finfo))))

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Environment file", "File size", "Modification time"})

	for _, fi := range finfo {
		t.AppendRow([]interface{}{hftx.Green(fi.Name()), hftx.Green(hf.SI(uint64(fi.Size()))),
			hftx.Green(fmt.Sprintf("%v", fi.ModTime().Format("2006/01/02 15:04:05")))})
	}
	t.SortBy([]table.SortBy{
		{Name: "Environment file", Mode: table.Asc},
		{Name: "File size", Mode: table.Asc},
	})
	t.SetStyle(table.StyleBold)
	t.Style().Format.Header = text.FormatDefault
	t.Render()

	return nil
}
