package config

import (
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"

	logger "github.com/informatiqal/qlik-test-users-tickets/Logger"
	"github.com/pelletier/go-toml/v2"
)

type config struct {
	Server struct {
		Port                 int
		HttpsCertificatePath string
	}
	Qlik struct {
		Host             string
		CertificatesPath string
		UserId           string
		UserDirectory    string
		DomainName       string
		Ports            struct {
			Repository int
			Proxy      int
		}
	}
	TestUsers struct {
		UserDirectory string
	}
}

var GlobalConfig config
var QlikClient *http.Client

func NewConfig() {
	log := logger.Zero
	pwd, _ := os.Executable()
	dir := filepath.Dir(pwd)

	configContent, readError := os.ReadFile(dir + "/config.toml")
	if readError != nil {
		log.Fatal().Err(readError).Msg("")
	}

	parseError := toml.Unmarshal([]byte(configContent), &GlobalConfig)
	if parseError != nil {
		log.Fatal().Err(readError).Msg("")
	}

	if GlobalConfig.Server.HttpsCertificatePath == "" {
		log.Fatal().Err(readError).Msg("Certificate path should be provided")
	}

	setQlikHttpClient()
}

func setQlikHttpClient() {
	log := logger.Zero

	qlikCert, err := tls.LoadX509KeyPair(
		GlobalConfig.Qlik.CertificatesPath+"/client.pem",
		GlobalConfig.Qlik.CertificatesPath+"/client_key.pem",
	)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{qlikCert},
	}

	client := &http.Client{Transport: customTransport}

	QlikClient = client
}
