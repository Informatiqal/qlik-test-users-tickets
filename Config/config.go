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
	TrustAllCerts    bool
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

var QlikClients = make(map[string]clientDetails)

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
	log := logger.Zero

	for q := range GlobalConfig.Qlik {
		qlikCert, err := tls.LoadX509KeyPair(
			GlobalConfig.Qlik[q].CertificatesPath+"/client.pem",
			GlobalConfig.Qlik[q].CertificatesPath+"/client_key.pem",
		)
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}

		trustAllCerts := true
		if GlobalConfig.Qlik[q].TrustAllCerts {
			trustAllCerts = !GlobalConfig.Qlik[q].TrustAllCerts
		}

		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: trustAllCerts,
			Certificates:       []tls.Certificate{qlikCert},
		}

		client := &http.Client{Transport: customTransport}

		details := clientDetails{
			UserId:        "",
			UserDirectory: "",
			HTTP:          client,
		}

		// if userId OR userDirectory are not provided use the default INTERNAL\sa_api
		if GlobalConfig.Qlik[q].UserId == "" || GlobalConfig.Qlik[q].UserDirectory == "" {
			details.UserId = "sa_api"
			details.UserDirectory = "INTERNAL"
		} else {
			details.UserId = GlobalConfig.Qlik[q].UserId
			details.UserDirectory = GlobalConfig.Qlik[q].UserDirectory
		}

		QlikClients[q] = details
	}
}
