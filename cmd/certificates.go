package cmd

import (
	"log"
	"os"

	util "github.com/informatiqal/qlik-test-users-tickets/Util"
	"github.com/spf13/cobra"
)

var certificatesCmd = &cobra.Command{
	Use:     "gencert",
	Short:   "Generate set of self-signed pem certificates in the current folder",
	Long:    "Generate set of self-signed pem certificates in the current folder",
	Example: ".\\qlik-test-user-tickets gencert",
	Run: func(_ *cobra.Command, args []string) {
		util.CreateSelfSignedCertificates()

		log.Print("Certificates generated")
		os.Exit(0)
	},
}
