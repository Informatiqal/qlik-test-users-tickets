package cmd

import (
	// "fmt"

	"log"
	"os"

	// logger "github.com/informatiqal/qlik-test-users-tickets/Logger"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

// type program struct{}

// var log = logger.Zero

var svcConfig = &service.Config{
	Name:        "QlikSenseTestUsers",
	DisplayName: "Qlik Sense Test Users",
	Description: "Generate Qlik tickets for test users",
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Install/Uninstall the windows service",
	Long:  "Install/Uninstall the windows service",
	Annotations: map[string]string{
		"command_category": "sub",
	},
}

var installServiceCmd = (&cobra.Command{
	Use:     "install",
	Args:    cobra.ExactArgs(0),
	Short:   "Install Qlik Test Users Tickets as a Windows service",
	Long:    "Install Qlik Test Users Tickets as a Windows service",
	Example: ".\\qlik-test-user-tickets service install",

	Run: func(_ *cobra.Command, args []string) {

		// prg := &program{}
		s, err := service.New(nil, svcConfig)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = s.Install()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Print("Service QlikSenseTestUsers installed")
		os.Exit(0)
	},
})

var uninstallServiceCmd = (&cobra.Command{
	Use:     "uninstall",
	Args:    cobra.ExactArgs(0),
	Short:   "Uninstall Qlik Test Users Tickets Windows service",
	Long:    "Uninstall Qlik Test Users Tickets Windows service",
	Example: ".\\qlik-test-user-tickets service uninstall",

	Run: func(_ *cobra.Command, args []string) {

		// prg := &program{}
		s, err := service.New(nil, svcConfig)
		if err != nil {
			log.Fatal(err.Error())
		}

		err = s.Uninstall()
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Print("Service QlikSenseTestUsers un-installed")
		os.Exit(0)
	},
})

func init() {
	serviceCmd.AddCommand(installServiceCmd, uninstallServiceCmd)
}
