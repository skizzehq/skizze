package counters

import (
	"counts/config"
	"counts/utils"
	"os"
	"path/filepath"
	"testing"
)

func setupTests() {
	os.Setenv("COUNTS_DATA_DIR", "/tmp/count_data")
	os.Setenv("COUNTS_INFO_DIR", "/tmp/count_info")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	os.Setenv("COUNTS_CONFIG", configPath)
}

func TestNoCounters(t *testing.T) {
	setupTests()
	var manager = GetManager()
	domains, err := manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if len(domains) != 0 {
		t.Error("Expected 0 counters, got", len(domains))
	}
}

func TestImmutableCounter(t *testing.T) {
	setupTests()
	var manager = GetManager()
	err := manager.CreateDomain("marvel", "immutable", 10000000)
	if err != nil {
		t.Error("Expected no errors while creating domain, got", err)
	}

	domains, err := manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors while getting domains, got", err)
	}
	if len(domains) != 1 {
		t.Error("Expected 1 counters, got", len(domains))
	}

	err = manager.AddToDomain("marvel", []string{"hulk", "thor"})
	if err != nil {
		t.Error("Expected no errors while adding to domain, got", err)
	}

	count, err := manager.GetCountForDomain("marvel")
	if len(domains) != 1 {
		t.Error("Expected 1 counters, got", len(domains))
	}

	if count != 2 {
		t.Error("Expected count == 2, got", count)
	}

	err = manager.DeleteDomain("marvel")
	if err != nil {
		t.Error("Expected no errors while deleting domain, got", err)
	}

	domains, err = manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors while getting domains, got", err)
	}
	if len(domains) != 0 {
		t.Error("Expected 0 counters, got", len(domains))
	}

}

func TestMutableCounter(t *testing.T) {
	setupTests()
	manager := GetManager()
	err := manager.CreateDomain("marvel", "mutable", 10000000)
	if err != nil {
		t.Error("Expected no errors while creating domain, got", err)
	}

	domains, err := manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors while getting domains, got", err)
	}
	if len(domains) != 1 {
		t.Error("Expected 1 counters, got", len(domains))
	}

	err = manager.AddToDomain("marvel", []string{"hulk", "thor"})
	if err != nil {
		t.Error("Expected no errors while adding to domain, got", err)
	}

	count, err := manager.GetCountForDomain("marvel")
	if len(domains) != 1 {
		t.Error("Expected 1 counters, got", len(domains))
	}

	if count != 2 {
		t.Error("Expected count == 2, got", count)
	}

	err = manager.DeleteDomain("marvel")
	if err != nil {
		t.Error("Expected no errors while deleting domain, got", err)
	}

	domains, err = manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors while getting domains, got", err)
	}
	if len(domains) != 0 {
		t.Error("Expected 0 counters, got", len(domains))
	}
}

func TestLoadInfo(t *testing.T) {
	setupTests()
	conf := config.GetConfig()
	testFilePath := filepath.Join(conf.GetInfoDir(), "test.json")

	var exists bool
	m := newManager()
	if _, exists = m.info["test"]; exists {
		t.Fatalf("expected test to not be initially loaded by manager")
	}

	f, err := os.Create(testFilePath)
	defer os.Remove(testFilePath)
	if err != nil {
		t.Fatal("Couldn't create test file")
	}
	f.WriteString(`{
		"id": "test",
		"type": "immutable",
		"capacity": 12345
	}`)

	m = newManager()
	if _, exists = m.info["test"]; !exists {
		t.Fatalf("expected test to be in loaded by manager")
	}
	// info := m.info["test"]

}
