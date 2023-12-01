//go:generate goversioninfo

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/informatiqal/qlik-test-users-tickets/Config"
	"github.com/informatiqal/qlik-test-users-tickets/Logger"
	"github.com/informatiqal/qlik-test-users-tickets/api"
	"github.com/informatiqal/qlik-test-users-tickets/cmd"
	"github.com/informatiqal/qlik-test-users-tickets/static"
	"github.com/kardianos/service"
)

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
	h2 := logger.Chain.Then(http.HandlerFunc(api.ProxyServiceList))
	h3 := logger.Chain.Then(http.HandlerFunc(api.TestUsersList))

	fs := http.FileServer(frontend.BuildHTTPFS())

	http.Handle("/", fs)
	http.Handle("/healthcheck", h)
	http.Handle("/api/ticket", h1)
	http.Handle("/api/proxies", h2)
	http.Handle("/api/users", h3)

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

	var version = "0.3.0"

	// execute and terminal command first
	cmd.Execute(version)

	// no arguments were provided - run the main logic
	// this is here mostly because of the --help command
	// have to check if it can be overwritten so the app to be
	// in control of the exit
	if len(os.Args) == 1 {

		// this is mentioned multiple times. Have to be moved somewhere
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

		err = s.Run()
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
	}

}
