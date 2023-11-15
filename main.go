package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

func main() {
	var testUsersDirectoryArg string
	var testUsersArg string
	var testUsers []string
	var certPathArg string
	var hostArg string

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

	flag.Parse()

	testUsers = strings.Split(testUsersArg, ";")

	if len(testUsers) == 1 && testUsers[0] == "" && testUsersDirectoryArg != "" {
		log.Fatal("User directory was provided but no users were passed")
	}

	if len(testUsers) > 0 &&
		testUsers[0] != "" &&
		testUsersDirectoryArg != "" &&
		certPathArg != "" &&
		hostArg != "" {
		createTestUsers(
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

func generateXrfkey() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func createTestUsers(host string, userDirectory string, users []string, certPath string) bool {

	cert, err := tls.LoadX509KeyPair(certPath+"/client.pem", certPath+"/client_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
	}

	client := &http.Client{Transport: customTransport}

	for _, user := range users {
		xrfkey := generateXrfkey()
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
		resp, err := client.Do(req)

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
