package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/informatiqal/qlik-test-users-tickets/API"
	"github.com/informatiqal/qlik-test-users-tickets/API/qlik"
	"github.com/informatiqal/qlik-test-users-tickets/Config"
	"github.com/informatiqal/qlik-test-users-tickets/Logger"
	"github.com/informatiqal/qlik-test-users-tickets/Util"
	"github.com/kardianos/service"
)

// var serviceLogger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Should be non-blocking, so run async using goroutine
	go p.run()
	return nil
}

func (p *program) run() {
	log := logger.Zero

	config.NewConfig()

	h := logger.Chain.Then(http.HandlerFunc(api.HealthCheckHandler))
	h1 := logger.Chain.Then(http.HandlerFunc(api.GenerateTicket))
	h2 := logger.Chain.Then(http.HandlerFunc(api.VirtualProxiesList))
	h3 := logger.Chain.Then(http.HandlerFunc(api.TestUsersList))

	http.Handle("/healthcheck", h)
	http.Handle("/ticket", h1)
	http.Handle("/virtualproxies", h2)
	http.Handle("/users", h3)
	// http.HandleFunc("/temp/", api.Test)

	log.Info().
		Msg("HTTPS server starting listening on port " + fmt.Sprint(config.GlobalConfig.Server.Port))
	err := http.ListenAndServeTLS(
		":"+fmt.Sprint(config.GlobalConfig.Server.Port),
		config.GlobalConfig.Server.HttpsCertificatePath+"/cert.pem",
		config.GlobalConfig.Server.HttpsCertificatePath+"/cert_key.pem",
		nil,
	)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
}

func (p *program) Stop(s service.Service) error {
	// Should be non-blocking
	return nil
}

func main() {
	log := logger.Zero

	var testUsersDirectoryArg string
	var testUsersArg string
	var testUsers []string
	var certPathArg string
	var hostArg string
	var generateCert bool
	var mode string

	svcConfig := &service.Config{
		Name:        "QlikSenseTestUsers",
		DisplayName: "Qlik Sense Test Users",
		Description: "Generate Qlik tickets for test users",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	flag.StringVar(&mode, "mode", "", "install/uninstall/run")

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

	if mode == "install" {
		err = s.Install()
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
	}

	if mode == "uninstall" {
		err = s.Uninstall()
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
	}

	if mode == "" || mode == "run" {
		err = s.Run()
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
	}

	if generateCert {
		util.CreateSelfSignedCertificates()
		os.Exit(0)
	}

	testUsers = strings.Split(testUsersArg, ";")

	if len(testUsers) == 1 && testUsers[0] == "" && testUsersDirectoryArg != "" {
		err := errors.New("user directory was provided but no users were passed")
		log.Fatal().Err(err).Msg("User directory was provided but no users were passed")
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

}
