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

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

func setupTests() {
	utils.PanicOnError(os.Setenv("SKZ_DATA_DIR", "/tmp/skizze_data"))
	utils.PanicOnError(os.Setenv("SKZ_INFO_DIR", "/tmp/skizze_info"))
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	utils.PanicOnError(os.Setenv("SKZ_CONFIG", configPath))
	tearDownTests()
}

func tearDownTests() {
	utils.PanicOnError(os.RemoveAll(config.GetConfig().DataDir))
	utils.PanicOnError(os.RemoveAll(config.GetConfig().InfoDir))
	utils.PanicOnError(os.Mkdir(config.GetConfig().DataDir, 0777))
	utils.PanicOnError(os.Mkdir(config.GetConfig().InfoDir, 0777))
	utils.PanicOnError(storage.CloseInfoDB())
	sketchesManager.Destroy()
}

func httpRequest(s *Server, t *testing.T, method string, sketch string, body string) *httptest.ResponseRecorder {
	reqBody := strings.NewReader(body)
	fullSketch := "http://skizze.io/" + sketch
	req, err := http.NewRequest(method, fullSketch, reqBody)
	if err != nil {
		t.Fatalf("%s", err)
	}
	respw := httptest.NewRecorder()
	s.ServeHTTP(respw, req)
	return respw
}

func unmarshalSketchsResult(resp *httptest.ResponseRecorder) sketchesResult {
	body, _ := ioutil.ReadAll(resp.Body)
	var r sketchesResult
	utils.PanicOnError(json.Unmarshal(body, &r))
	return r
}

func unmarshalSketchResult(resp *httptest.ResponseRecorder) sketchResult {
	body, _ := ioutil.ReadAll(resp.Body)
	var r sketchResult
	utils.PanicOnError(json.Unmarshal(body, &r))
	return r
}

func TestSketchsInitiallyEmpty(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s, err := New()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	resp := httpRequest(s, t, "GET", "", "")
	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}
	result := unmarshalSketchsResult(resp)
	if len(result.Result) != 0 {
		t.Fatalf("Initial resultCount != 0. Got %d", len(result.Result))
	}
}

func TestPost(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s, err := New()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	resp := httpRequest(s, t, "POST", "hllpp/marvel", `{
		"capacity": 100000
	}`)

	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}

	resp = httpRequest(s, t, "GET", "", `{}`)
	result := unmarshalSketchsResult(resp)
	if len(result.Result) != 1 {
		t.Fatalf("after add resultCount != 1. Got %d", len(result.Result))
	}
}

func TestHLL(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s, err := New()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	resp := httpRequest(s, t, "POST", "hllpp/marvel", `{
		"capacity": 100000
	}`)

	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}

	resp = httpRequest(s, t, "GET", "", `{}`)
	result := unmarshalSketchsResult(resp)
	if len(result.Result) != 1 {
		t.Fatalf("after add resultCount != 1. Got %d", len(result.Result))
	}

	resp = httpRequest(s, t, "PUT", "hllpp/marvel", `{
		"values": ["magneto", "wasp", "beast"]
	}`)

	resp = httpRequest(s, t, "GET", "hllpp/marvel", `{}`)
	result2 := unmarshalSketchResult(resp)

	if result2.Result.(float64) != 3 {
		t.Fatalf("after add resultCount != 1. Got %f.0", result2.Result.(float64))
	}
}

func TestCML(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s, err := New()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	resp := httpRequest(s, t, "POST", "cml/x-force", `{
		"epsilon": 0.05, "delta": 0.99
	}`)

	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}

	resp = httpRequest(s, t, "GET", "", `{}`)
	result := unmarshalSketchsResult(resp)
	if len(result.Result) != 1 {
		t.Fatalf("after add resultCount != 1. Got %d", len(result.Result))
	}

	resp = httpRequest(s, t, "PUT", "cml/x-force", `{
		"values": ["magneto", "wasp", "beast", "magneto"]
	}`)

	resp = httpRequest(s, t, "GET", "cml/x-force", `{"values": ["magneto"]}`)
	result2 := unmarshalSketchResult(resp).Result.(map[string]interface{})

	if v, ok := result2["magneto"]; ok && uint(v.(float64)) != 2 {
		t.Fatalf("after add resultCount != 2. Got %d", uint(v.(float64)))
	}
}

func TestTopK(t *testing.T) {
	setupTests()
	defer tearDownTests()
	s, err := New()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	resp := httpRequest(s, t, "POST", "topk/x-force", `{
			"capacity": 3
	}`)

	if resp.Code != 200 {
		t.Fatalf("Invalid Response Code %d - %s", resp.Code, resp.Body.String())
		return
	}

	resp = httpRequest(s, t, "GET", "", `{}`)
	result := unmarshalSketchsResult(resp)
	if len(result.Result) != 1 {
		t.Fatalf("after add resultCount != 1. Got %d", len(result.Result))
	}

	resp = httpRequest(s, t, "PUT", "topk/x-force", `{
			"values": ["magneto", "wasp", "beast", "magneto", "pyro"]
		}`)

	resp = httpRequest(s, t, "GET", "topk/x-force", `{"values":[]}`)

	result2 := unmarshalSketchResult(resp).Result.([]interface{})
	res := make([]map[string]interface{}, len(result2))
	for i, v := range result2 {
		res[i] = v.(map[string]interface{})
	}

	if v, ok := res[0]["Key"]; ok && v.(string) != "magneto" {
		t.Fatalf("Expected \"magneto\" in first position, got, %s", v.(string))
	}
}
