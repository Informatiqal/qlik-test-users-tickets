package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/informatiqal/qlik-test-users-tickets/API"
	"github.com/informatiqal/qlik-test-users-tickets/API/qlik"
	"github.com/informatiqal/qlik-test-users-tickets/Config"
	"github.com/informatiqal/qlik-test-users-tickets/Util"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func main() {
	appLogFile, _ := os.OpenFile(
		"app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	appLog := zerolog.MultiLevelWriter(os.Stdout, appLogFile)
	logger := zerolog.New(appLog).With().Timestamp().Logger()

	httpLogFile, _ := os.OpenFile(
		"http.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	httpLog := zerolog.MultiLevelWriter(httpLogFile)
	httpLogger := zerolog.New(httpLog).With().Timestamp().Logger()

	c := alice.New()
	c = c.Append(hlog.NewHandler(httpLogger))

	c = c.Append(
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}),
	)
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

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

	// initialize the config (aka read the config file)
	config.NewConfig()

	h := c.Then(http.HandlerFunc(api.HealthCheckHandler))
	h1 := c.Then(http.HandlerFunc(api.GenerateTicket))
	h2 := c.Then(http.HandlerFunc(api.VirtualProxiesList))
	h3 := c.Then(http.HandlerFunc(api.TestUsersList))

	http.Handle("/healthcheck", h)
	http.Handle("/ticket", h1)
	http.Handle("/virtualproxies", h2)
	http.Handle("/users", h3)
	// http.HandleFunc("/temp/", api.Test)

	logger.Info().
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
