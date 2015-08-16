package counters

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/seiflotfy/counts/config"
	"github.com/seiflotfy/counts/counters/abstract"
	"github.com/seiflotfy/counts/storage"
	"github.com/seiflotfy/counts/utils"
)

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
	storage.CloseInfoDB()
}

func TestNoCounters(t *testing.T) {
	setupTests()
	defer tearDownTests()
	var manager, err = GetManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	domains, err := manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if len(domains) != 0 {
		t.Error("Expected 0 counters, got", len(domains))
	}
}

func TestDefaultCounter(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var manager, err = GetManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	err = manager.CreateDomain("marvel", "default", 10000000)
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

func TestPurgableCounter(t *testing.T) {
	setupTests()
	defer tearDownTests()

	manager, err := GetManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	err = manager.CreateDomain("marvel", "purgable", 10000000)
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

func TestDumpLoadDefaultInfo(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var exists bool
	m1, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists = m1.info["avengers"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	err = m1.CreateDomain("avengers", "default", 1000000)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists = m2.info["avengers"]; !exists {
		t.Error("expected avengers to be in loaded by manager")
	}
}

func TestDumpLoadDefaultData(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var exists bool
	m1, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists = m1.info["avengers"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	m1.CreateDomain("avengers", "default", 1000000)

	m1.AddToDomain("avengers", []string{"sabertooth",
		"thunderbolt", "havoc", "cyclops"})

	res, err := m1.GetCountForDomain("avengers")
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	if res != 4 {
		t.Error("expected avengers to have count 4, got", res)
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	res, err = m2.GetCountForDomain("avengers")
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	if res != 4 {
		t.Error("expected avengers to have count 4, got", res)
	}
}

func TestDumpLoadPurgableInfo(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var exists bool
	m1, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists = m1.info["avengers"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	err = m1.CreateDomain("avengers", abstract.Purgable, 1000000)
	if err != nil {
		t.Fatal(err)
	}

	err = m1.AddToDomain("avengers", []string{"hulk", "storm"})
	if err != nil {
		t.Fatal(err)
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists = m2.info["avengers"]; !exists {
		t.Error("expected avengers to be in loaded by manager")
	}

	count, err := m2.GetCountForDomain("avengers")
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	if count != 2 {
		t.Error("Expected count == 2, got", count)
	}

}

func TestExtremeParallelDefaultCounter(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists := m1.info["avengers"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	m1.CreateDomain("avengers", "default", 1000000)
	m1.CreateDomain("x-men", "default", 1000000)

	fd, err := os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	i := 0
	values := []string{} //{"a", "aam"} //"doorknob", "doorless"
	for scanner.Scan() {
		s := []byte(scanner.Text())
		values = append(values, string(s))
		i++
	}

	// Add all values in a go routine per value
	var wg sync.WaitGroup
	defer wg.Wait()
	resChan := make(chan interface{})

	var pFunc = func(value string) {
		defer wg.Done()
		m1.AddToDomain("avengers", []string{value})
		resChan <- nil
	}
	for _, value := range values {
		wg.Add(1)
		go pFunc(value)
	}
	for j := 0; j < len(values); j++ {
		<-resChan
	}

	// add all values in one bulk
	m1.AddToDomain("x-men", values)
	count1, err := m1.GetCountForDomain("avengers")
	count2, err := m1.GetCountForDomain("x-men")
	if count1 != count2 {
		t.Error("expected avengers count == x-men count, got", count1, "!=", count2)
	}
}

func TestFailCreateDomain(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}

	// test for unknown domainType
	err = m1.CreateDomain("marvel", "wrong", 10000000)
	if err == nil {
		t.Error("Expected errors while creating domain, got", err)
	}

	buffer := make([]byte, config.MaxKeySize+1)
	for i := 0; i < config.MaxKeySize+1; i++ {
		buffer[i] = byte(49) // ascii 1
	}
	domainID := string(buffer)
	// test for too long domainID
	err = m1.CreateDomain(domainID, "default", 10000000)
	if err == nil {
		t.Error("Expected errors while creating domain, got", err)
	}
}

func TestFailDeleteDomain(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}
	err = m1.DeleteDomain("-1")
	if err == nil {
		t.Error("Expected error, got", err)
	}
}

func TestFailDeleteFromDomain(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}
	err = m1.DeleteFromDomain("-1", []string{})
	if err == nil {
		t.Error("Expected error, got", err)
	}
}

func TestFailGetCountForDomain(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}
	_, err = m1.GetCountForDomain("-1")
	if err == nil {
		t.Error("Expected error, got", err)
	}
}
