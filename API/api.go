package api

import (
	"encoding/json"
	"fmt"
	"github.com/informatiqal/qlik-test-users-tickets/API/qlik"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GenerateTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	qlikTicket, err := qlik.CreateTicketForUser()
	// qlikTicket, err := qlik.GetVirtualProxies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type ReqBody struct {
		User string `json:"user"`
	}
	var reqBody ReqBody

	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type TicketResponse struct {
		Ticket string `json:"ticket"`
	}

	var response TicketResponse
	response.Ticket = qlikTicket

	b, err := json.Marshal(&response)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err)
	}
}

func VirtualProxiesList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	virtualProxies, err := qlik.GetVirtualProxies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(&virtualProxies)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err)
	}
}
