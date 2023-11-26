package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/informatiqal/qlik-test-users-tickets/API/qlik"
	"github.com/rs/zerolog/log"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	log.Debug().Msg("Healthcheck api called")

	w.WriteHeader(http.StatusOK)
}

func GenerateTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	type CreateTicketParams struct {
		User               string        `json:"userId"`
		VirtualProxyPrefix string        `json:"virtualProxyPrefix"`
		Attributes         []interface{} `json:"attributes"`
		ProxyId            string        `json:"proxyId"`
	}

	var reqBody CreateTicketParams
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Malformed JSON", http.StatusBadRequest)
		return
	}

	if reqBody.User == "" {
		http.Error(w, "\"user\" value is mandatory", http.StatusBadRequest)
		return
	}

	if reqBody.ProxyId == "" {
		http.Error(w, "\"proxyId\" value is mandatory", http.StatusBadRequest)
		return
	}

	isTestUser, userDetails := validateUser(reqBody.User)
	if !isTestUser {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	attributes, err := json.Marshal(reqBody.Attributes)
	if err != nil {
		log.Debug().Err(err).Msg("")
	}

	qlikTicket, err := qlik.CreateTicketForUser(
		reqBody.User,
		userDetails.UserDirectory,
		reqBody.VirtualProxyPrefix,
		string(attributes),
		reqBody.ProxyId,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}

	b, err := json.Marshal(&qlikTicket)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err)
	}
}

func ProxyServiceList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	proxyServices, err := qlik.GetProxyServices()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(&proxyServices)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err)
	}
}

func TestUsersList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	users, err := qlik.GetTestUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(&users)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err)
	}
}

func validateUser(userId string) (bool, *qlik.User) {
	user, err := qlik.GetUserDetails(userId)
	if err != nil {
		return false, nil
	}

	isTestUserDirectory := strings.HasPrefix(user.UserDirectory, "TESTING_")

	return isTestUserDirectory, user
}
