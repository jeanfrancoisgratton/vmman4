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
	env.Port = uint(hf.GetIntValFromPrompt("Please enter the db port (empty means 5432): "))
	if env.Port > 65535 || env.Port == 0 {
		env.Port = 5432
	}
	env.DbName = hf.GetStringValFromPrompt("Enter the database name (empty means vmmandb): ")
	if env.DbName == "" {
		env.DbName = "vmmandb"
	}
	env.User = hf.GetStringValFromPrompt("Enter the db username: ")
	env.Password = hf.EncodeString(hf.GetPassword("Enter the username's password: ", shared.DebugMode), "")

	return env.SaveEnvironmentFile()
}

// Save the above structure into a JSON file in the user's config directory

func (e DBenvironmentStruct) SaveEnvironmentFile() *ce.CustomError {
	jStream, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return &ce.CustomError{Title: "Unable to marshal data", Message: err.Error()}
	}
	if err = os.WriteFile(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", "dbinfo.json"),
		jStream, 0600); err != nil {
		return &ce.CustomError{Title: err.Error()}
	}

	hftx.EnabledSign("Database environment file created")
	return nil
}

// Load the JSON environment file in the user's db connection config, and store it into a data type (struct)

func LoadEnvironmentFile() (DBenvironmentStruct, *ce.CustomError) {
	var payload DBenvironmentStruct
	var err error

	jFile, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".config", "JFG", "vmman4", "dbinfo.json"))
	if err != nil {
		return DBenvironmentStruct{}, &ce.CustomError{Title: err.Error()}
	}
	err = json.Unmarshal(jFile, &payload)
	if err != nil {
		return DBenvironmentStruct{}, &ce.CustomError{Title: err.Error()}
	} else {
		payload.Password = hf.DecodeString(payload.Password, "")
		return payload, nil
	}
}
