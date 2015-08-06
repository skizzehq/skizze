package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type requestData struct {
	domain     string
	domainType string
	values     []string
}

type server struct {
}

func newServer() *server {
	server := server{}
	return &server
}

func (svr *server) handleTopRequest(w http.ResponseWriter, method string, data requestData) {
	switch {
	case method == "GET":
		// Get all counters
		break
	case method == "POST":
		// Create new counter
		break
	case method == "DELETE":
		// Remove values from domain
		break
	}
	fmt.Fprintf(w, "Huh?")
}

func (svr *server) handleDomainRequest(w http.ResponseWriter, method string, data requestData) {
	switch {
	case method == "GET":
		// Get a count for a specific domain
		break
	case method == "POST":
		// Add values to domain
		break
	case method == "DELETE":
		// Delete Counter
		break
	}
	fmt.Fprintf(w, "Huh?")
}

func (svr *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.URL.Path[1:])
	method := r.Method
	body, _ := ioutil.ReadAll(r.Body)

	var data requestData
	_ = json.Unmarshal(body, &data)

	if len(domain) == 0 {
		svr.handleTopRequest(w, method, data)
	} else {
		data.domain = domain
		svr.handleDomainRequest(w, method, data)
	}
}

/*
Run ...
*/
func (svr *server) Run() {
	http.ListenAndServe(":8080", svr)
}

/*
Stop ...
*/
func (svr *server) Stop() {
	os.Exit(0)
}
