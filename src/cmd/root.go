// vmman4
// src/cmd/root.go

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vmman4",
	Short:   "Virtual Machines Manager",
	Version: "0.100-0 (2024.07.28)",
	Long: `With this tool you will be able to :
- Start/stop/create/destroy/edit virtual machines
- Manage the hypervisors where those VMs run
- Manage a PGSQL database back-end that will store VM metadata`,
}

var changelogCmd = &cobra.Command{
	Use:     "changelog",
	Aliases: []string{"cl"},
	Short:   "Application changelog",
	Run: func(cmd *cobra.Command, args []string) {
		changeLog()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.DisableAutoGenTag = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(changelogCmd)
}

func changeLog() {
	fmt.Print(`VERSION			DATE			COMMENT
-------			----			-------
0.100			2024.07.28		initial version
`)
}
