package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = ""

var rootCmd = &cobra.Command{
	Use:  "qs_test_users",
	Long: `Generate Qlik Sense (Windows) tickets for test users`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// func Execute(mainVersion, branchName, commitSha string) {
func Execute(mainVersion string) {
	version = mainVersion
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints an error message so we just want to exit
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(certificatesCmd)
	rootCmd.AddCommand(usersCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(ticketCmd)
}
