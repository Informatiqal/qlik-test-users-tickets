package qlik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/informatiqal/qlik-test-users-tickets/Config"
	util "github.com/informatiqal/qlik-test-users-tickets/Util"
	"log"
	"net/http"
	"strings"
)

type VirtualProxy struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Prefix      string `json:"prefix"`
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
			log.Fatal(err)
		}

		req.Header.Add("X-Qlik-Xrfkey", xrfkey)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Qlik-User", "UserDirectory=INTERNAL;UserId=sa_api")
		resp, err := config.QlikClient.Do(req)

		if err != nil {
			log.Fatal(err)
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

func CreateTicketForUser() (string, error) {
	return "abc", nil
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
		log.Fatal(err)
	}

	req.Header.Add("X-Qlik-Xrfkey", xrfkey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Qlik-User", "UserDirectory=INTERNAL;UserId=sa_api")
	resp, err := config.QlikClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var responseData []VirtualProxy

	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&responseData)

	return &responseData, nil
}
