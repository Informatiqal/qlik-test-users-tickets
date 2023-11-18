package qlik

import (
	"bytes"
	"encoding/json"
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
	Id          string `json:"id"`
	Description string `json:"description"`
	Prefix      string `json:"prefix"`
}

type User struct {
	Id            string `json:"id"`
	UserDirectory string `json:"userDirectory"`
	Name          string `json:"name"`
}

type GeneratedTicket struct {
	UserId        string `json:"userId"`
	UserDirectory string `json:"userDirectory"`
	Ticket        string `json:"ticket"`
}

var log zerolog.Logger

func init() {
	log = logger.Zero
}

func CreateTestUsers(host string, userDirectory string, users []string, certPath string) bool {
	for _, user := range users {
		xrfkey := util.GenerateXrfkey()
		url := fmt.Sprintf("https://%s:4242/qrs/user?Xrfkey=%s", host, xrfkey)

		jsonBody := []byte(
			fmt.Sprintf(
				`{"userId": "%s","userDirectory": "%s","removedExternally": false,"blacklisted": false,"name": "%s"}`,
				strings.TrimSpace(user),
				userDirectory,
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
		req.Header.Add("X-Qlik-User", "UserDirectory=INTERNAL;UserId=sa_api")
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

func CreateTicketForUser(
	userId string,
	userDirectory string,
	vp string,
) (GeneratedTicket, error) {
	var vpString string
	if vp != "" {
		vpString = vp + "/"
	} else {
		vpString = ""
	}

	xrfkey := util.GenerateXrfkey()
	url := fmt.Sprintf(
		"https://%s:4243/qps/%sticket?Xrfkey=%s",
		config.GlobalConfig.Qlik.Host,
		vpString,
		xrfkey,
	)

	jsonBody := []byte(
		fmt.Sprintf(
			`{"userId": "%s","userDirectory": "%s"}`,
			strings.TrimSpace(userId),
			userDirectory,
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
	req.Header.Add("X-Qlik-User", "UserDirectory=INTERNAL;UserId=sa_api")
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := GeneratedTicket{}
		return t, err
	}

	var responseData GeneratedTicket

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	msg := fmt.Sprintf(
		"Ticket %s was generated for userId %s",
		responseData.Ticket,
		userId,
	)
	log.Info().Msg(msg)

	return responseData, nil
}

func GetVirtualProxies() (*[]VirtualProxy, error) {
	xrfkey := util.GenerateXrfkey()
	url := fmt.Sprintf(
		"https://%s:4242/qrs/virtualproxyconfig?Xrfkey=%s",
		config.GlobalConfig.Qlik.Host,
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
	req.Header.Add("X-Qlik-User", "UserDirectory=INTERNAL;UserId=sa_api")
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

func GetTestUsers() (*[]User, error) {
	xrfkey := util.GenerateXrfkey()

	baseUrl := "https://" + config.GlobalConfig.Qlik.Host + ":4242/qrs/user"

	encoded := url.Values{}
	encoded.Set("Xrfkey", xrfkey)
	encoded.Set("filter", "userDirectory sw 'TEST'")

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
	req.Header.Add("X-Qlik-User", "UserDirectory=INTERNAL;UserId=sa_api")
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

	baseUrl := "https://" + config.GlobalConfig.Qlik.Host + ":4242/qrs/user"
	baseUrl += "/" + userId

	encoded := url.Values{}
	encoded.Set("Xrfkey", xrfkey)

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
	req.Header.Add("X-Qlik-User", "UserDirectory=INTERNAL;UserId=sa_api")
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("")
		t := User{}
		return &t, err
	}

	var responseData User

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	return &responseData, nil
}
