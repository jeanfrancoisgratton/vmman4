// vmman4
// src/cmd/root.go

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
	Use:   "vmman",
	Short: "Virtual Machine Management Tool",
	Long: `This tool allows you to manage your VM farm.
With it you can start, stop, snapshot, snapshot-revert, create or delete VMs.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.White("vmman4 v0.30.0 (2026.09.09), Go version = v" + strings.TrimPrefix(runtime.Version(), "go")))
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

	//rootCmd.AddCommand(clCmd)
	rootCmd.PersistentFlags().BoolVarP(&shared.QuietOutput, "quiet", "q", false, "Suppress output")
	rootCmd.PersistentFlags().BoolVarP(&shared.DebugMode, "debug", "D", false, "Enable debug mode")
	rootCmd.PersistentFlags().StringVarP(&shared.ConnectionFilename, "connectionfile", "c", "", "Connection configuration file")
	rootCmd.PersistentFlags().StringVarP(&shared.ConnectURI, "connectionuri", "C", "qemu:///system", "Connection URI")
}
