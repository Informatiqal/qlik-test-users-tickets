package config

import (
	"crypto/tls"
	"net/http"
	"os"
	"path/filepath"

	logger "github.com/informatiqal/qlik-test-users-tickets/Logger"
	"github.com/pelletier/go-toml/v2"
)

type qlikCluster struct {
	CertificatesPath string
	UserId           string
	UserDirectory    string
	RepositoryHost   string
	DomainMapping    map[string]string
}

type config struct {
	Server struct {
		Port                 int
		HttpsCertificatePath string
	}
	Qlik map[string]qlikCluster
}

var GlobalConfig config

type clientDetails struct {
	UserId        string
	UserDirectory string
	HTTP          *http.Client
}

var QlikClients map[string]clientDetails

func NewConfig() {
	log := logger.Zero
	pwd, _ := os.Executable()
	dir := filepath.Dir(pwd)

	configContent, readError := os.ReadFile(dir + "/config.toml")
	if readError != nil {
		log.Fatal().Err(readError).Msg(readError.Error())
	}

	parseError := toml.Unmarshal([]byte(configContent), &GlobalConfig)
	if parseError != nil {
		log.Fatal().Err(readError).Msg(parseError.Error())
	}

	if GlobalConfig.Server.Port == 0 {
		GlobalConfig.Server.Port = 8081
		log.Warn().Msg("Port value was not provided in the config! Using default 8081")
	}

	if GlobalConfig.Server.HttpsCertificatePath == "" {
		log.Fatal().Err(readError).Msg("Certificate path should be provided")
	}

	setQlikHttpClients()
}

func setQlikHttpClients() {

	for q := range GlobalConfig.Qlik {
		log := logger.Zero

		qlikCert, err := tls.LoadX509KeyPair(
			GlobalConfig.Qlik[q].CertificatesPath+"/client.pem",
			GlobalConfig.Qlik[q].CertificatesPath+"/client_key.pem",
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

		details := clientDetails{
			UserId:        "",
			UserDirectory: "",
			HTTP:          client,
		}

		// if userId is not provided use the default INTERNAL\sa_api
		if GlobalConfig.Qlik[q].UserId == "" {
			details.UserId = "sa_api"
			details.UserDirectory = "INTERNAL"
		}

		QlikClients[q] = details
	}
}
