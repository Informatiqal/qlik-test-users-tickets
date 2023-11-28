package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "qlik-test-user-tickets",
	Long: `Generate Qlik Sense (Windows) tickets for test users`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// func Execute(mainVersion, branchName, commitSha string) {
func Execute() {
	// version = mainVersion
	// branch = branchName
	// commit = commitSha
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints an error message so we just want to exit
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(certificatesCmd)
	rootCmd.AddCommand(usersCmd)
}
