// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/dbCommands.go
// Original timestamp : 2026.10.02 14:34:31

package cmd

import (
	"fmt"
	"vmman4/db_mgt"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:     "db",
	Aliases: []string{"database"},
	Example: "vmman db init",
	Short:   "Database support subcommands",
	Long:    `You need to provide one of the subcommands: init, resync, status, inventory.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You need to provide one of the following subcommands: init, resync, status, inventory")
	},
}

var dbEnvCmd = &cobra.Command{
	Use:     "env",
	Example: "vmman db env",
	Short:   "Creates the database environment file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db_mgt.CreateDBenvFile(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

var dbInitCmd = &cobra.Command{
	Use: "init",
	//Aliases: []string{"ls"},
	Example: "vmman db init",
	Short:   "Creates the database",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db_mgt.DBinit(); err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)

	dbCmd.AddCommand(dbInitCmd, dbEnvCmd)
}
