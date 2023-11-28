package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Version will output the current build information",
	Long:    "Version will output the current build information",
	Example: ".\\qs_test_users version",

	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("version: %s", version)
		os.Exit(0)
	},
}
