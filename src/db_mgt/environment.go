// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/db_mgt/environment.go
// Original timestamp : 2026.10.02 15:18:10

package db_mgt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"vmman4/shared"

	ce "github.com/jeanfrancoisgratton/customError/v3"
	hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// CreateDBenvFile: the user will be prompted for the db credentials
// The user that you will provide must have ownership of the database (that you will also provide)

func CreateDBenvFile() *ce.CustomError {
	var env DBenvironmentStruct
	//fmt.Println("The root dir value should be an absolute path, and all other values relative to it")

	hftx.WarningSign("IMPORTANT")
	fmt.Println("The database must already exist, as well the user. Moreover, that user must have ownership of the DB")

	env.Host = hf.GetStringValFromPrompt("Enter the database hostname: ")
	env.Port = uint(hf.GetIntValFromPrompt("Please enter the db port: "))
	if env.Port > 65535 || env.Port == 0 {
		env.Port = 5432
	}
	env.User = hf.GetStringValFromPrompt("Enter the db username: ")
	env.Password = hf.EncodeString(hf.GetPassword("Enter the username's password: ", shared.DebugMode), "")

	jStream, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return &ce.CustomError{Title: "Unable to marshal data", Message: err.Error()}
	}
	rcFile := filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", "dbenv.json")
	if err = os.WriteFile(rcFile, jStream, 0600); err != nil {
		return &ce.CustomError{Title: "Unable to write file", Message: err.Error()}
	}

	hftx.EnabledSign("Database environment file created")

	return nil
}
