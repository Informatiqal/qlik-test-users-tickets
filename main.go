package main

import (
	"flag"
	"fmt"
	"github.com/informatiqal/qlik-test-users-tickets/API"
	"github.com/informatiqal/qlik-test-users-tickets/API/qlik"
	"github.com/informatiqal/qlik-test-users-tickets/Config"
	"github.com/informatiqal/qlik-test-users-tickets/Util"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var testUsersDirectoryArg string
	var testUsersArg string
	var testUsers []string
	var certPathArg string
	var hostArg string
	var generateCert bool

	flag.StringVar(
		&testUsersDirectoryArg,
		"userDirectory",
		"",
		"Test users user directory (default is TESTING)",
	)

	flag.StringVar(
		&hostArg,
		"host",
		"",
		"Qlik Sense host/machine name. Make sure that the machine is accessible on port 4242",
	)

	flag.StringVar(
		&certPathArg,
		"certPath",
		"",
		"Directory location where cert.pem and cert_key.pem are located",
	)

	flag.StringVar(
		&testUsersArg,
		"users",
		"",
		"Semicolon list of users to be created into to the provided user directory",
	)

	flag.BoolVar(
		&generateCert,
		"generateCert",
		false,
		"Generate self-signed certificates in the current folder",
	)

	flag.Parse()

	if generateCert {
		util.CreateSelfSignedCertificates()
		os.Exit(0)
	}

	testUsers = strings.Split(testUsersArg, ";")

	if len(testUsers) == 1 && testUsers[0] == "" && testUsersDirectoryArg != "" {
		log.Fatal("User directory was provided but no users were passed")
	}

	if len(testUsers) > 0 &&
		testUsers[0] != "" &&
		testUsersDirectoryArg != "" &&
		certPathArg != "" &&
		hostArg != "" {
		qlik.CreateTestUsers(
			hostArg,
			testUsersDirectoryArg,
			testUsers,
			certPathArg,
		)

		fmt.Println("")
		fmt.Println("Operation completed!")
		defer os.Exit(0)
	}

	// initialize the config (aka read the config file)
	config.NewConfig()

	http.HandleFunc("/healthcheck", api.HealthCheckHandler)
	http.HandleFunc("/ticket", api.GenerateTicket)
	http.HandleFunc("/virtualproxies", api.VirtualProxiesList)

	log.Printf("HTTPS server starting listening on port %v\n", config.GlobalConfig.Server.Port)
	err := http.ListenAndServeTLS(
		":"+fmt.Sprint(config.GlobalConfig.Server.Port),
		config.GlobalConfig.Server.HttpsCertificatePath+"/cert.pem",
		config.GlobalConfig.Server.HttpsCertificatePath+"/cert_key.pem",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

}
