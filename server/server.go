package server

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

/*
Server managed the http connections and communciates with the counters manager
*/
type Server struct {
}

/*
New returns a new Server
*/
func New() *Server {
	server := Server{}
	return &server
}

func (srv *Server) handleTopRequest(w http.ResponseWriter, method string, data requestData) {
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

func (srv *Server) handleDomainRequest(w http.ResponseWriter, method string, data requestData) {
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

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.URL.Path[1:])
	method := r.Method
	body, _ := ioutil.ReadAll(r.Body)

	var data requestData
	_ = json.Unmarshal(body, &data)

	if len(domain) == 0 {
		srv.handleTopRequest(w, method, data)
	} else {
		data.domain = domain
		srv.handleDomainRequest(w, method, data)
	}
}

/*
Run ...
*/
func (srv *Server) Run() {
	http.ListenAndServe(":8080", srv)
}

/*
Stop ...
*/
func (srv *Server) Stop() {
	os.Exit(0)
}
