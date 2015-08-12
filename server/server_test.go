package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/seiflotfy/counts/config"
	"github.com/seiflotfy/counts/utils"
)

type domainsResult struct {
	Result []string `json:"result"`
	Error  error    `json:"error"`
}

type domainResult struct {
	Result uint  `json:"result"`
	Error  error `json:"error"`
}

func setupTests() {
	os.Setenv("COUNTS_DATA_DIR", "/tmp/count_data")
	os.Setenv("COUNTS_INFO_DIR", "/tmp/count_info")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	os.Setenv("COUNTS_CONFIG", configPath)
	tearDownTests()
}

func tearDownTests() {
	os.RemoveAll(config.GetConfig().GetDataDir())
	os.RemoveAll(config.GetConfig().GetInfoDir())
	os.Mkdir(config.GetConfig().GetDataDir(), 0777)
	os.Mkdir(config.GetConfig().GetInfoDir(), 0777)
}

func request(s *Server, t *testing.T, method string, domain string, body string) *httptest.ResponseRecorder {
	reqBody := strings.NewReader(body)
	fullDomain := "http://counters.io/" + domain
	req, err := http.NewRequest(method, fullDomain, reqBody)
	if err != nil {
		t.Fatalf("%s", err)
	}
	respw := httptest.NewRecorder()
	s.ServeHTTP(respw, req)
	return respw
}

func unmarshalDomainsResult(resp *httptest.ResponseRecorder) domainsResult {
	body, _ := ioutil.ReadAll(resp.Body)
	var r domainsResult
	json.Unmarshal(body, &r)
	return r
}

func unmarshalDomainResult(resp *httptest.ResponseRecorder) domainResult {
	body, _ := ioutil.ReadAll(resp.Body)
	var r domainResult
	json.Unmarshal(body, &r)
	return r
}

func TestDomainsInitiallyEmpty(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s := New()
	resp := request(s, t, "GET", "", "")
	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}
	result := unmarshalDomainsResult(resp)
	if len(result.Result) != 0 {
		t.Fatalf("Initial resultCount != 0. Got %s", result)
	}
}

func TestCreateDomain(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s := New()
	resp := request(s, t, "POST", "marvel", `{
		"domainType": "immutable",
		"capacity": 100000
	}`)

	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}

	resp = request(s, t, "GET", "", `{}`)
	result := unmarshalDomainsResult(resp)
	if len(result.Result) != 1 {
		t.Fatalf("after add resultCount != 1. Got %s", result)
	}
}

func TestHLL(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s := New()
	resp := request(s, t, "POST", "marvel", `{
		"domainType": "immutable",
		"capacity": 100000
	}`)

	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}

	resp = request(s, t, "GET", "", `{}`)
	result := unmarshalDomainsResult(resp)
	if len(result.Result) != 1 {
		t.Fatalf("after add resultCount != 1. Got %s", result)
	}

	resp = request(s, t, "PUT", "marvel", `{
		"values": ["magneto", "wasp", "beast"]
	}`)

	resp = request(s, t, "GET", "marvel", `{}`)
	result2 := unmarshalDomainResult(resp)
	if result2.Result != 3 {
		t.Fatalf("after add resultCount != 1. Got %d", result2.Result)
	}
}
