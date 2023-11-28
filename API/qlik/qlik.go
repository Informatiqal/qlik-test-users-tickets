package qlik

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	sysLog "log"
	"net/http"
	"net/url"
	"strings"

	"github.com/informatiqal/qlik-test-users-tickets/Config"
	"github.com/informatiqal/qlik-test-users-tickets/Logger"
	util "github.com/informatiqal/qlik-test-users-tickets/Util"
	"github.com/rs/zerolog"
)

type VirtualProxy struct {
	// Id          string `json:"id"`
	Description string `json:"description"`
	Prefix      string `json:"prefix"`
}

type ProxyServiceRaw struct {
	Id                      string `json:"id"`
	ServerNodeConfiguration struct {
		Name     string `json:"name"`
		HostName string `json:"hostName"`
	} `json:"serverNodeConfiguration"`
	Settings struct {
		VirtualProxies []VirtualProxy `json:"virtualProxies"`
	} `json:"settings"`
}

type ProxyService struct {
	Id             string         `json:"id"`
	Name           string         `json:"name"`
	HostName       string         `json:"hostName"`
	VirtualProxies []VirtualProxy `json:"virtualProxies"`
}

type User struct {
	UserId        string `json:"userId"`
	UserDirectory string `json:"userDirectory"`
	Name          string `json:"name"`
}

type GeneratedTicket struct {
	UserId             string `json:"userId"`
	UserDirectory      string `json:"userDirectory"`
	Ticket             string `json:"ticket"`
	VirtualProxyPrefix string `json:"virtualProxyPrefix"`
	Links              struct {
		Qmc string `json:"qmc"`
		Hub string `json:"hub"`
	} `json:"links"`
}

var log zerolog.Logger

func init() {
	log = logger.Zero
}

func CreateTestUsers(
	host string,
	userDirectorySuffix string,
	users []string,
	certPath string,
) bool {
	for _, user := range users {
		xrfkey := util.GenerateXrfkey()
		url := fmt.Sprintf("https://%s:4242/qrs/user?Xrfkey=%s", host, xrfkey)

		jsonBody := []byte(
			fmt.Sprintf(
				`{"userId": "%s","userDirectory": "%s","removedExternally": false,"blacklisted": false,"name": "%s"}`,
				strings.TrimSpace(user),
				"TESTING_"+userDirectorySuffix,
				strings.TrimSpace(user),
			),
		)
		bodyReader := bytes.NewReader(jsonBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bodyReader,
		)
		if err != nil {
			sysLog.Fatal(err.Error())
		}

		req.Header.Add("X-Qlik-Xrfkey", xrfkey)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Qlik-User", fmt.Sprintf(
			"UserDirectory=%s;UserId=%s",
			config.GlobalConfig.Qlik.UserDirectory,
			config.GlobalConfig.Qlik.UserId,
		))
		resp, err := config.QlikClient.Do(req)

		if err != nil {
			sysLog.Fatal(err.Error())
		}

		type userDetails struct {
			Id string `json:"id"`
		}

		var responseData userDetails

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&responseData)

		fmt.Printf("User \"%s\" created -> %s\n", strings.TrimSpace(user), responseData.Id)
	}

	return true
}

// TODO: this should return something meaningful and not boolean
func CreateTestUsersCmd(
	users []string,
	userDirectorySuffix string,
) bool {
	for _, user := range users {
		xrfkey := util.GenerateXrfkey()
		url := fmt.Sprintf(
			"https://%s:4242/qrs/user?Xrfkey=%s",
			config.GlobalConfig.Qlik.RepositoryHost,
			xrfkey,
		)

		jsonBody := []byte(
			fmt.Sprintf(
				`{"userId": "%s","userDirectory": "%s","removedExternally": false,"blacklisted": false,"name": "%s"}`,
				strings.TrimSpace(user),
				"TESTING_"+userDirectorySuffix,
				strings.TrimSpace(user),
			),
		)
		bodyReader := bytes.NewReader(jsonBody)

		req, err := http.NewRequest(
			"POST",
			url,
			bodyReader,
		)
		if err != nil {
			sysLog.Fatal(err.Error())
		}

		req.Header.Add("X-Qlik-Xrfkey", xrfkey)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add(
			"X-Qlik-User",
			fmt.Sprintf(
				"UserDirectory=%s;UserId=%s",
				config.GlobalConfig.Qlik.UserDirectory,
				config.GlobalConfig.Qlik.UserId,
			),
		)
		resp, err := config.QlikClient.Do(req)

		if err != nil {
			sysLog.Fatal(err.Error())
		}

		type userDetails struct {
			Id string `json:"id"`
		}

		var responseData userDetails

		decoder := json.NewDecoder(resp.Body)
		decoder.Decode(&responseData)

		fmt.Printf("User \"%s\" created -> %s\n", strings.TrimSpace(user), responseData.Id)
	}

	return true
}

// TODO: this should just create ticket and nothing more (refactor)
func CreateTicketForUser(
	userId string,
	userDirectory string,
	vp string,
	attributes string,
	proxyId string,
) (GeneratedTicket, error) {
	var vpString string
	if vp != "" {
		vpString = vp + "/"
	} else {
		vpString = ""
	}

	proxyService, err := getProxyService(proxyId)
	if err != nil {
		log.Error().Err(err).Msg("Proxy service id not found")
		t := GeneratedTicket{}
		return t, err
	}

	xrfkey := util.GenerateXrfkey()
	url := fmt.Sprintf(
		"https://%s:4243/qps/%sticket?Xrfkey=%s",
		proxyService.ServerNodeConfiguration.HostName,
		vpString,
		xrfkey,
	)

	jsonBody := []byte(
		fmt.Sprintf(
			`{"userId": "%s","userDirectory": "%s", "attributes": %s}`,
			strings.TrimSpace(userId),
			userDirectory,
			attributes,
		),
	)

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(
		"POST",
		url,
		bodyReader,
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := GeneratedTicket{}
		return t, err
	}

	req.Header.Add("X-Qlik-Xrfkey", xrfkey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Qlik-User", fmt.Sprintf(
		"UserDirectory=%s;UserId=%s",
		config.GlobalConfig.Qlik.UserDirectory,
		config.GlobalConfig.Qlik.UserId,
	))
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := GeneratedTicket{}
		return t, err
	}

	if resp.StatusCode != http.StatusCreated {
		log.Error().
			Err(err).
			Msg("Error while generating the ticket. Proxy API responded with: " + resp.Status)
		t := GeneratedTicket{}
		return t, errors.New("")
	}

	var responseData GeneratedTicket

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	presentationUrl := getPresentationUrl(proxyService.ServerNodeConfiguration.HostName)

	responseData.VirtualProxyPrefix = vp
	responseData.Links.Qmc = "https://" + presentationUrl + "/" + vpString + "qmc?qlikTicket=" + responseData.Ticket
	responseData.Links.Hub = "https://" + presentationUrl + "/" + vpString + "hub?qlikTicket=" + responseData.Ticket

	msg := fmt.Sprintf(
		`Ticket "%s" was generated for userId "%s" in virtual proxy "%s" on proxy node "%s"`,
		responseData.Ticket,
		userId,
		vp,
		proxyService.ServerNodeConfiguration.HostName,
	)
	log.Info().Msg(msg)

	return responseData, nil
}

func GetVirtualProxies() (*[]VirtualProxy, error) {
	xrfkey := util.GenerateXrfkey()
	url := fmt.Sprintf(
		"https://%s:4242/qrs/virtualproxyconfig?Xrfkey=%s",
		config.GlobalConfig.Qlik.RepositoryHost,
		xrfkey,
	)

	req, err := http.NewRequest(
		"GET",
		url,
		http.NoBody,
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := []VirtualProxy{}
		return &t, err
	}

	req.Header.Add("X-Qlik-Xrfkey", xrfkey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Qlik-User", fmt.Sprintf(
		"UserDirectory=%s;UserId=%s",
		config.GlobalConfig.Qlik.UserDirectory,
		config.GlobalConfig.Qlik.UserId,
	))
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := []VirtualProxy{}
		return &t, err
	}

	var responseData []VirtualProxy

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	return &responseData, nil
}

func GetProxyServices() (*[]ProxyService, error) {
	xrfkey := util.GenerateXrfkey()
	baseUrl := fmt.Sprintf(
		"https://%s:4242/qrs/ProxyService/full",
		config.GlobalConfig.Qlik.RepositoryHost,
	)

	encoded := url.Values{}
	encoded.Set("Xrfkey", xrfkey)
	encoded.Set("filter", "(serverNodeConfiguration.proxyEnabled eq True)")

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?%s", baseUrl, encoded.Encode()),
		http.NoBody,
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := []ProxyService{}
		return &t, err
	}

	req.Header.Add("X-Qlik-Xrfkey", xrfkey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Qlik-User", fmt.Sprintf(
		"UserDirectory=%s;UserId=%s",
		config.GlobalConfig.Qlik.UserDirectory,
		config.GlobalConfig.Qlik.UserId,
	))
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := []ProxyService{}
		return &t, err
	}

	var responseData []ProxyServiceRaw

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	var apiData []ProxyService = []ProxyService{}

	for i := 0; i < len(responseData); i++ {
		p := ProxyService{
			Id:             responseData[i].Id,
			Name:           responseData[i].ServerNodeConfiguration.Name,
			HostName:       responseData[i].ServerNodeConfiguration.HostName,
			VirtualProxies: responseData[i].Settings.VirtualProxies,
		}

		apiData = append(apiData, p)
	}

	return &apiData, nil
}

func GetTestUsers() (*[]User, error) {
	xrfkey := util.GenerateXrfkey()

	baseUrl := "https://" + config.GlobalConfig.Qlik.RepositoryHost + ":4242/qrs/user"

	encoded := url.Values{}
	encoded.Set("Xrfkey", xrfkey)
	encoded.Set("filter", "userDirectory sw 'TESTING_'")

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?%s", baseUrl, encoded.Encode()),
		http.NoBody,
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := []User{}
		return &t, err
	}

	req.Header.Add("X-Qlik-Xrfkey", xrfkey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Qlik-User", fmt.Sprintf(
		"UserDirectory=%s;UserId=%s",
		config.GlobalConfig.Qlik.UserDirectory,
		config.GlobalConfig.Qlik.UserId,
	))
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := []User{}
		return &t, err
	}

	var responseData []User

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	return &responseData, nil
}

func GetUserDetails(userId string) (*User, error) {
	xrfkey := util.GenerateXrfkey()

	baseUrl := "https://" + config.GlobalConfig.Qlik.RepositoryHost + ":4242/qrs/user"

	encoded := url.Values{}
	encoded.Set("Xrfkey", xrfkey)
	encoded.Set("filter", "(userId eq '"+userId+"')")

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s?%s", baseUrl, encoded.Encode()),
		http.NoBody,
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := User{}
		return &t, err
	}

	req.Header.Add("X-Qlik-Xrfkey", xrfkey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Qlik-User", fmt.Sprintf(
		"UserDirectory=%s;UserId=%s",
		config.GlobalConfig.Qlik.UserDirectory,
		config.GlobalConfig.Qlik.UserId,
	))
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := User{}
		return &t, err
	}

	var responseData []User

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)
	if len(responseData) == 0 {
		return nil, errors.New("User not found")
	}

	return &responseData[0], nil
}

func getProxyService(id string) (*ProxyServiceRaw, error) {
	xrfkey := util.GenerateXrfkey()
	url := fmt.Sprintf(
		"https://%s:4242/qrs/proxyservice/%s?Xrfkey=%s",
		config.GlobalConfig.Qlik.RepositoryHost,
		id,
		xrfkey,
	)

	req, err := http.NewRequest(
		"GET",
		url,
		http.NoBody,
	)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := ProxyServiceRaw{}
		return &t, err
	}

	req.Header.Add("X-Qlik-Xrfkey", xrfkey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Qlik-User", fmt.Sprintf(
		"UserDirectory=%s;UserId=%s",
		config.GlobalConfig.Qlik.UserDirectory,
		config.GlobalConfig.Qlik.UserId,
	))
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := ProxyServiceRaw{}
		return &t, err
	}

	if resp.StatusCode == 404 {
		log.Error().Err(err).Msg("ProxyService not found!" + id)
		t := ProxyServiceRaw{}
		return &t, errors.New("ProxyService not found")
	}

	var responseData ProxyServiceRaw

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	return &responseData, nil
}

func getPresentationUrl(machineName string) string {
	prettyName := config.GlobalConfig.Qlik.DomainMapping[machineName]

	if prettyName == "" {
		return machineName
	}

	return prettyName
}
