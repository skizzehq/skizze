package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/sketches"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

type requestData struct {
	id         string
	typ        string
	Properties map[string]float64 `json:"properties"`
	Values     []string           `json:"values"`
}

var logger = utils.GetLogger()
var sketchesManager *sketches.ManagerStruct

/*
Server manages the http connections and communciates with the sketches manager
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
	sketchesManager, err = sketches.GetManager()
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
		// Get all sketches
		sketches, err = sketchesManager.GetSketches()
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
		if _, err := w.Write(js); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (srv *Server) handleSketchRequest(w http.ResponseWriter, method string, data requestData) {
	var res sketchResult
	var err error

	// TODO (mb): handle errors from sketchesManager.*
	switch {
	case method == "GET":
		// Get a count for a specific sketch
		count, err := sketchesManager.GetCountForSketch(data.id, data.typ, data.Values)
		logger.Info.Printf("[%v]: Getting state for sketch: %v of type %s", method, data.id, data.typ)
		res = sketchResult{count, err}
	case method == "POST":
		// Create a new sketch counter
		err = sketchesManager.CreateSketch(data.id, data.typ, data.Properties)
		logger.Info.Printf("[%v]: Creating new sketch: %v of type %s", method, data.id, data.typ)
		res = sketchResult{0, err}
	case method == "PUT":
		// Add values to counter
		err = sketchesManager.AddToSketch(data.id, data.typ, data.Values)
		logger.Info.Printf("[%v]: Adding values to sketch: %v of type %s", method, data.id, data.typ)
		res = sketchResult{nil, err}
	case method == "PURGE":
		// Purges values from counter
		err = sketchesManager.DeleteFromSketch(data.id, data.typ, data.Values)
		logger.Info.Printf("[%v]: Purging values from sketch: %v of type %s", method, data.id, data.typ)
		res = sketchResult{nil, err}
	case method == "DELETE":
		// Delete Counter
		err := sketchesManager.DeleteSketch(data.id, data.typ)
		logger.Info.Printf("[%v]: Deleting sketch: %v of type %s", method, data.id, data.typ)
		res = sketchResult{nil, err}
	default:
		logger.Error.Printf("[%v]: Invalid Method: %v", method, http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("Invalid Method: %s", method), http.StatusBadRequest)
		return
	}

	if res.Error != nil {
		http.Error(w, fmt.Sprintf("Error with operation %s on %s: %s", method, data.id, res.Error.Error()), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(js); err != nil {
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

	if data.Properties == nil {
		data.Properties = make(map[string]float64)
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
	port := int(conf.Port)
	logger.Info.Println("Server up and running on port: " + strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), srv)
	utils.PanicOnError(err)
	err = gracehttp.Serve(&http.Server{Addr: ":" + strconv.Itoa(port), Handler: srv})
	utils.PanicOnError(err)
}

/*
Stop ...
*/
func (srv *Server) Stop() {
	//FIXME make sure everything is written to disk
	logger.Info.Println("Stopping server...")
	err := storage.CloseInfoDB()
	utils.PanicOnError(err)
	os.Exit(0)
}
