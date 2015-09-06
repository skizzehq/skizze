package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/counters"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

type requestData struct {
	id       string
	typ      string
	Capacity uint64   `json:"capacity"`
	Error    float64  `json:"error"`
	Storage  string   `json:"storage"`
	Values   []string `json:"values"`
}

var logger = utils.GetLogger()
var counterManager *counters.ManagerStruct

/*
Server manages the http connections and communciates with the counters manager
*/
type Server struct{}

type sketchesResult struct {
	Result []string `json:"result"`
	Error  error    `json:"error"`
}

type sketchResult struct {
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
	var sketches []string
	var js []byte

	switch {
	case method == "GET":
		// Get all counters
		sketches, err = counterManager.GetSketches()
		js, err = json.Marshal(sketchesResult{sketches, err})
		logger.Info.Printf("[%v]: Getting all available sketches", method)
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

func (srv *Server) handleSketchRequest(w http.ResponseWriter, method string, data requestData) {
	var res sketchResult
	var err error

	// TODO (mb): handle errors from counterManager.*
	switch {
	case method == "GET":
		// Get a count for a specific sketch
		count, err := counterManager.GetCountForSketch(data.id, data.typ, data.Values)
		logger.Info.Printf("[%v]: Getting counter for sketch: %v", method, data.id)
		res = sketchResult{count, err}
	case method == "POST":
		// Create a new sketch counter
		err = counterManager.CreateSketch(data.id, data.typ, data.Capacity)
		logger.Info.Printf("[%v]: Creating new sketch: %v", method, data.id)
		res = sketchResult{0, err}
	case method == "PUT":
		// Add values to counter
		err = counterManager.AddToSketch(data.id, data.typ, data.Values)
		logger.Info.Printf("[%v]: Updating counter for sketch: %v", method, data.id)
		res = sketchResult{nil, err}
	case method == "PURGE":
		// Purges values from counter
		err = counterManager.DeleteFromSketch(data.id, data.typ, data.Values)
		logger.Info.Printf("[%v]: Purging values for sketch: %v", method, data.id)
		res = sketchResult{nil, err}
	case method == "DELETE":
		// Delete Counter
		err := counterManager.DeleteSketch(data.id, data.typ)
		logger.Info.Printf("[%v]: Deleting sketch: %v", method, data.id)
		res = sketchResult{nil, err}
	default:
		logger.Error.Printf("[%v]: Invalid Method: %v", method, http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Invalid Method: %s", method), http.StatusBadRequest)
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
	method := r.Method
	paths := strings.Split(r.URL.Path[1:], "/")
	body, _ := ioutil.ReadAll(r.Body)
	var data requestData
	if len(body) > 0 {
		err := json.Unmarshal(body, &data)
		if err != nil {
			logger.Error.Printf("An error has ocurred: %v", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		data = requestData{}
	}

	if len(paths) == 1 {
		srv.handleTopRequest(w, method, data)
	} else if len(paths) == 2 {
		data.typ = strings.TrimSpace(string(paths[0]))
		data.id = strings.TrimSpace(strings.Join(paths[1:], "/"))
		srv.handleSketchRequest(w, method, data)
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
	logger.Info.Println("Stopping server...")
	storage.CloseInfoDB()
	os.Exit(0)
}
