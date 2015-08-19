package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/seiflotfy/counts/config"
	"github.com/seiflotfy/counts/counters"
	"github.com/seiflotfy/counts/storage"
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
type Server struct{}

type domainsResult struct {
	Result []string `json:"result"`
	Error  error    `json:"error"`
}

type domainResult struct {
	Result interface{} `json:"result"`
	Error  error       `json:"error"`
}

/*
New returns a new Server
*/
func New() (*Server, error) {
	var err error
	counterManager, err = counters.GetManager()
	if err != nil {
		return nil, err
	}
	server := Server{}
	return &server, nil
}

func (srv *Server) handleTopRequest(w http.ResponseWriter, method string, data requestData) {
	var err error
	var domains []string
	var js []byte

	switch {
	case method == "GET":
		// Get all counters
		domains, err = counterManager.GetDomains()
		js, err = json.Marshal(domainsResult{domains, err})
	case method == "MERGE":
		// Reserved for merging hyper log log
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
		return
	default:
		http.Error(w, "Invalid Method: "+method, http.StatusBadRequest)
		return
	}

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (srv *Server) handleDomainRequest(w http.ResponseWriter, method string, data requestData) {
	var res domainResult
	var err error

	// TODO (mb): handle errors from counterManager.*
	switch {
	case method == "GET":
		// Get a count for a specific domain
		count, err := counterManager.GetCountForDomain(data.Domain)
		res = domainResult{count, err}
	case method == "POST":
		// Create a new domain counter
		err = counterManager.CreateDomain(data.Domain, data.DomainType, data.Capacity)
		res = domainResult{0, err}
	case method == "PUT":
		// Add values to counter
		err = counterManager.AddToDomain(data.Domain, data.Values)
		res = domainResult{nil, err}
	case method == "PURGE":
		// Purges values from counter
		err = counterManager.DeleteFromDomain(data.Domain, data.Values)
		res = domainResult{nil, err}
	case method == "DELETE":
		// Delete Counter
		err := counterManager.DeleteFromDomain(data.Domain, data.Values)
		res = domainResult{nil, err}
	default:
		http.Error(w, "Invalid Method: "+method, http.StatusBadRequest)
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
	port := int(conf.GetPort())
	logger.Info.Println("Server up and running on port: " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), srv)
}

/*
Stop ...
*/
func (srv *Server) Stop() {
	storage.CloseInfoDB()
	os.Exit(0)
}
