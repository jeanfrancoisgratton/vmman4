// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/cmd/root.go
// Original timestamp : 2026.10.14 22:52:57

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"vmman4/shared"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vmman",
	Example: "vmman vm list\nvmman -c myhypervisor vm list",
	Short:   "Virtual Machine Management Tool",
	Long: `This tool allows you to manage your VM farm.
With it you can start, stop, snapshot, snapshot-revert, create or delete VMs.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// An explicitly supplied --db/-d flag takes precedence over the environment.
		// BoolVarP has already set shared.DBsupport from the flag when it was supplied.
		if cmd.Root().PersistentFlags().Changed("db") {
			return nil
		}

		// If the flag was not supplied, use VMMAN_DBSUPPORT when it is present.
		value, exists := os.LookupEnv("VMMAN_DBSUPPORT")
		if !exists {
			shared.DBsupport = false
			return nil
		}

		switch strings.ToLower(strings.TrimSpace(value)) {
		case "1", "true", "t":
			shared.DBsupport = true
		case "0", "false", "f":
			shared.DBsupport = false
		default:
			return fmt.Errorf("invalid VMMAN_DBSUPPORT value %q: expected 0, 1, t, f, true, or false", value)
		}

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Example: "vmman version",
	Short:   "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.White("vmman v" + buildVersion + " (" + buildDate + "), Go version = v" + strings.TrimPrefix(runtime.Version(), "go") + " (" + runtime.GOARCH + ")"))
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

	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().BoolVarP(&shared.QuietOutput, "quiet", "q", false, "Suppress output")
	rootCmd.PersistentFlags().BoolVarP(&shared.DebugMode, "debug", "D", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&shared.ConnectionFilename, "connectionfile", "c", "", "Connection configuration file")
	rootCmd.PersistentFlags().StringVarP(&shared.ConnectURI, "connectionuri", "C", "qemu:///system", "Connection URI")
}
