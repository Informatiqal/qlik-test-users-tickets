package config

import (
	"crypto/tls"
	"github.com/pelletier/go-toml/v2"
	"log"
	"net/http"
	"os"
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
	pwd, _ := os.Getwd()

	configContent, readError := os.ReadFile(pwd + "/config.toml")
	if readError != nil {
		log.Fatal(readError)
	}

	parseError := toml.Unmarshal([]byte(configContent), &GlobalConfig)
	if parseError != nil {
		panic(parseError)
	}

	if GlobalConfig.Server.HttpsCertificatePath == "" {
		log.Fatal("Certificate path should be provided")
	}

	setQlikHttpClient()
}

func setQlikHttpClient() {
	qlikCert, err := tls.LoadX509KeyPair(
		GlobalConfig.Qlik.CertificatesPath+"/client.pem",
		GlobalConfig.Qlik.CertificatesPath+"/client_key.pem",
	)
	if err != nil {
		log.Fatal(err)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{qlikCert},
	}

	client := &http.Client{Transport: customTransport}

	QlikClient = client
}
