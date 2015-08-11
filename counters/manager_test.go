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

func tearDownTests() {
	os.RemoveAll(config.GetConfig().GetDataDir())
	os.RemoveAll(config.GetConfig().GetInfoDir())
	os.Mkdir(config.GetConfig().GetDataDir(), 0777)
	os.Mkdir(config.GetConfig().GetInfoDir(), 0777)
}

func TestNoCounters(t *testing.T) {
	setupTests()
	defer tearDownTests()
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
	defer tearDownTests()

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
	defer tearDownTests()

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

func TestDumpLoadInfo(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var exists bool
	m1 := newManager()
	if _, exists = m1.info["avengers"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	m1.CreateDomain("avengers", "immutable", 1000000)

	m2 := newManager()
	if _, exists = m2.info["avengers"]; !exists {
		t.Error("expected avengers to be in loaded by manager")
	}
}

func TestDumpLoadImmutableData(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var exists bool
	m1 := newManager()
	if _, exists = m1.info["avengers"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	m1.CreateDomain("avengers", "immutable", 1000000)
	m1.AddToDomain("avengers", []string{"sabertooth", "thunderbolt", "havoc", "cyclops"})

	res, err := m1.GetCountForDomain("avengers")
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	if res != 4 {
		t.Error("expected avengers to have count 4, got", res)
	}

	m2 := newManager()
	res, err = m2.GetCountForDomain("avengers")
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	if res != 4 {
		t.Error("expected avengers to have count 4, got", res)
	}

}
