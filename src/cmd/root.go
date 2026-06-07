// vmman4
// src/cmd/root.go

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"vmman4/shared"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "vmman4",
	Short:   "Virtual Machine Management Tool",
	Version: "DEBUG0.10.00 (2026.05.18)",
	Long: `This tool allows you to manage your VM farm.
With it you can start, stop, snapshot, snapshot-revert, create or delete VMs.`,
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

	//rootCmd.AddCommand(clCmd)
	rootCmd.PersistentFlags().BoolVarP(&shared.QuietOutput, "quiet", "q", false, "Suppress output")
	rootCmd.PersistentFlags().BoolVarP(&shared.DebugMode, "debug", "D", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&shared.ConnectionFilename, "connectionfile", "c", "", "Connection configuration file")
	rootCmd.PersistentFlags().StringVarP(&shared.ConnectURI, "connectionuri", "C", "qemu:///system", "Connection URI")
}
