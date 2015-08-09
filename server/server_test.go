package server

import (
	"counts/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type domainsResult struct {
	Result []string `json:"result"`
	Error  error    `json:"error"`
}

func setupTests() {
	os.Setenv("COUNTS_DATA_DIR", "/tmp/count_data")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "data/default_config.json")
	os.Setenv("COUNTS_CONFIG", configPath)
}

func request(s *Server, t *testing.T, method string, domain string, body string) *httptest.ResponseRecorder {
	reqBody := strings.NewReader(body)
	req, err := http.NewRequest(method, "http://counters.io/"+domain, reqBody)
	if err != nil {
		t.Fatalf("%s", err)
	}
	respw := httptest.NewRecorder()
	s.ServeHTTP(respw, req)
	return respw
}

func unmarschal(resp *httptest.ResponseRecorder) domainsResult {
	body, _ := ioutil.ReadAll(resp.Body)
	var r domainsResult
	json.Unmarshal(body, &r)
	return r
}

func TestDomainsInitiallyEmpty(t *testing.T) {
	setupTests()
	s := New()
	resp := request(s, t, "GET", "", "{}")
	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}
	result := unmarschal(resp)
	if len(result.Result) != 0 {
		t.Fatalf("Initial resultCount != 0. Got %s", result)
	}
}

func TestCreateDomain(t *testing.T) {
	setupTests()
	s := New()
	resp := request(s, t, "POST", "", `{
		"domain": "marvel",
		"domainType": "immutable",
		"capacity": 100000,
		"values": []
	}`)

	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}

	resp = request(s, t, "GET", "", `{}`)
	result := unmarschal(resp)
	if len(result.Result) != 1 {
		t.Fatalf("after add resultCount != 1. Got %s", result)
	}
}
