package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/seiflotfy/counts/config"
	"github.com/seiflotfy/counts/counters"
	"github.com/seiflotfy/counts/utils"
)

type requestData struct {
	Domain     string   `json:"domain"`
	DomainType string   `json:"domainType"`
	Capacity   uint64   `json:"capacity"`
	Values     []string `json:"values"`
}

var logger = utils.GetLogger()
var counterManager *counters.ManagerStruct

/*
Server manages the http connections and communciates with the counters manager
*/
type Server struct {
}

type result struct {
	Result interface{} `json:"result"`
	Error  error       `json:"error"`
}

/*
New returns a new Server
*/
func New() *Server {
	counterManager = counters.GetManager()
	server := Server{}
	return &server
}

func (srv *Server) handleTopRequest(w http.ResponseWriter, method string, data requestData) {
	var res result
	switch {
	case method == "GET":
		// Get all counters
		domains, err := counterManager.GetDomains()
		res = result{domains, err}
	case method == "MERGE":
		// Reserved for merging hyper log log
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
		return
	default:
		http.Error(w, "Invalid Method: "+method, http.StatusBadRequest)
		return
	}

	// Somebody tried a PUT request (ignore)
	if res.Result == nil && res.Error == nil {
		fmt.Fprintf(w, "Huh?")
		return
	}

	js, err := json.Marshal(res)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (srv *Server) handleDomainRequest(w http.ResponseWriter, method string, data requestData) {
	var res result
	switch {
	case method == "GET":
		// Get a count for a specific domain
		count, err := counterManager.GetCountForDomain(data.Domain)
		res = result{count, err}
	case method == "POST":
		// Create a new domain counter
		err := counterManager.CreateDomain(data.Domain, data.DomainType, data.Capacity)
		res = result{data.Domain, err}
	case method == "PUT":
		// Add values to counter
		err := counterManager.AddToDomain(data.Domain, data.Values)
		res = result{nil, err}
	case method == "PURGE":
		// Purges values from counter
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	case method == "DELETE":
		// Delete Counter
		err := counterManager.DeleteFromDomain(data.Domain, data.Values)
		res = result{nil, err}
	}
	// Somebody tried a PUT request (ignore)
	if res.Result == nil && res.Error == nil {
		logger.Error.Println(w, "Huh?")
		return
	}

	js, err := json.Marshal(res)

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domain := strings.TrimSpace(r.URL.Path[1:])
	method := r.Method
	body, _ := ioutil.ReadAll(r.Body)

	var data requestData
	if len(body) > 0 {
		err := json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		data = requestData{}
	}

	if len(domain) == 0 {
		srv.handleTopRequest(w, method, data)
	} else {
		data.Domain = domain
		srv.handleDomainRequest(w, method, data)
	}
}

/*
Run ...
*/
func (srv *Server) Run() {
	conf := config.GetConfig()
	port := conf.GetPort()
	logger.Info.Println("Server up and running on port: " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), srv)
}

/*
Stop ...
*/
func (srv *Server) Stop() {
	os.Exit(0)
}
