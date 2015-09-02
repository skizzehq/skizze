package counters

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/counters/abstract"
	"github.com/seiflotfy/skizze/counters/wrappers/topk"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func setupTests() {
	os.Setenv("SKZ_DATA_DIR", "/tmp/skizze_manager_data")
	os.Setenv("SKZ_INFO_DIR", "/tmp/skizze_manager_info")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	os.Setenv("SKZ_CONFIG", configPath)
	tearDownTests()
}

func tearDownTests() {
	storage.CloseInfoDB()
	os.RemoveAll(config.GetConfig().GetDataDir())
	os.RemoveAll(config.GetConfig().GetInfoDir())
	os.Mkdir(config.GetConfig().GetDataDir(), 0777)
	os.Mkdir(config.GetConfig().GetInfoDir(), 0777)
	manager.Destroy()
}

func TestNoCounters(t *testing.T) {
	setupTests()
	defer tearDownTests()
	var manager, err = newManager()
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

func TestDuplicateCounters(t *testing.T) {
	setupTests()
	defer tearDownTests()
	var manager, err = newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	err = manager.CreateDomain("marvel", "cardinality", 10000000)
	if err != nil {
		t.Error("Expected no errors while creating domain, got", err)
	}
	err = manager.CreateDomain("marvel", "topk", 10000000)
	if err == nil {
		t.Error("Expected errors while creating domain duplicate domain, got", err)
	}
}

func TestDefaultCounter(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var manager, err = newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	domains, err := manager.GetDomains()
	if err != nil {
		t.Error("Expected no errors while getting domains, got", err)
	}
	if len(domains) != 0 {
		t.Error("Expected 0 counters, got", len(domains))
	}

	err = manager.CreateDomain("marvel", "cardinality", 10000000)
	if err != nil {
		t.Error("Expected no errors while creating domain, got", err)
	}

	domains, err = manager.GetDomains()
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

	count, err := manager.GetCountForDomain("marvel", nil)
	if len(domains) != 1 {
		t.Error("Expected 1 counters, got", len(domains))
	}

	if count.(uint) != 2 {
		t.Error("Expected count == 2, got", count.(uint))
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

	manager, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	err = manager.CreateDomain("marvel", abstract.PurgableCardinality, 10000000)
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

	count, err := manager.GetCountForDomain("marvel", nil)
	if count.(uint) != 2 {
		t.Error("Expected count == 2, got", count.(uint))
	}

	err = manager.DeleteFromDomain("marvel", []string{"hulk"})
	if err != nil {
		t.Error("Expected no errors while getting domains, got", err)
	}

	count, err = manager.GetCountForDomain("marvel", nil)
	if count.(uint) != 1 {
		t.Error("Expected count == 1, got", count.(uint))
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
	if _, exists = m1.info["x-force"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	err = m1.CreateDomain("x-force", "cardinality", 1000000)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists = m2.info["x-force"]; !exists {
		t.Error("expected x-force to be in loaded by manager")
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
	m1.CreateDomain("avengers", "cardinality", 1000000)

	m1.AddToDomain("avengers", []string{"sabertooth",
		"thunderbolt", "havoc", "cyclops"})

	res, err := m1.GetCountForDomain("avengers", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	if res.(uint) != 4 {
		t.Error("expected avengers to have count 4, got", res.(uint))
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	res, err = m2.GetCountForDomain("avengers", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	if res.(uint) != 4 {
		t.Error("expected avengers to have count 4, got", res.(uint))
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
	err = m1.CreateDomain("avengers", abstract.PurgableCardinality, 1000000)
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

	count, err := m2.GetCountForDomain("avengers", nil)
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	if count.(uint) != 2 {
		t.Error("Expected count == 2, got", count.(uint))
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
	m1.CreateDomain("avengers", "cardinality", 1000000)
	m1.CreateDomain("x-men", "cardinality", 1000000)

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
		if i == 10000 {
			break
		}
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
	count1, err := m1.GetCountForDomain("avengers", nil)
	count2, err := m1.GetCountForDomain("x-men", nil)
	if count1 != count2 {
		t.Error("expected avengers count == x-men count, got", count1, "!=", count2)
	}
}

func TestExtremeParallelPurgableCounter(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists := m1.info["avengers"]; exists {
		t.Error("expected avengers to not be initially loaded by manager")
	}
	m1.CreateDomain("avengers", abstract.PurgableCardinality, 1000000)
	m1.CreateDomain("x-men", abstract.PurgableCardinality, 1000000)

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
		if i == 10000 {
			break
		}
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
	count1, err := m1.GetCountForDomain("avengers", nil)
	count2, err := m1.GetCountForDomain("x-men", nil)
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
	err = m1.CreateDomain(domainID, "cardinality", 10000000)
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
	_, err = m1.GetCountForDomain("-1", nil)
	if err == nil {
		t.Error("Expected error, got", err)
	}
}

func TestTopKCounter(t *testing.T) {
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
	m1.CreateDomain("avengers", "topk", 3)

	m1.AddToDomain("avengers", []string{"sabertooth",
		"thunderbolt", "havoc", "cyclops", "cyclops", "cyclops", "havoc"})

	res, err := m1.GetCountForDomain("avengers", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	top := res.([]topk.ResultElement)
	if len(top) != 3 {
		t.Error("expected avengers to have 3 elements, got", len(top))
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	res, err = m2.GetCountForDomain("avengers", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	top = res.([]topk.ResultElement)
	if len(top) != 3 {
		t.Error("expected avengers to have 3 elements, got", len(top))
	}

	if top[0].Key != "cyclops" {
		t.Error("expected 1st avengers key == cyclops, got", top[0].Key)
	}
	if top[0].Count != 3 {
		t.Error("expected 1st avengers count == 3, got", top[0].Count)
	}
	if top[1].Key != "havoc" {
		t.Error("expected 1st avengers key == havoc, got", top[1].Key)
	}
	if top[1].Count != 2 {
		t.Error("expected 1st avengers count == 2, got", top[1].Count)
	}
}

func TestCMLCounter(t *testing.T) {
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
	m1.CreateDomain("avengers", abstract.Frequency, 3)

	m1.AddToDomain("avengers", []string{"sabertooth",
		"thunderbolt", "havoc", "cyclops", "cyclops", "cyclops", "havoc"})

	res, err := m1.GetCountForDomain("avengers", []string{"cyclops"})
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	res, err = m2.GetCountForDomain("avengers", []string{"cyclops"})
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	counts := res.(map[string]uint)
	if v, ok := counts["cyclops"]; !ok {
		t.Error("expected to find 'cyclops' in avengers, got", ok)
	} else if v != 3 {
		t.Error("expected 'cyclops' count == 3, got", v)
	}

	res, err = m2.GetCountForDomain("avengers", []string{"havoc"})
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	counts = res.(map[string]uint)
	if v, ok := counts["havoc"]; !ok {
		t.Error("expected to find 'havoc' in avengers, got", ok)
	} else if v != 2 {
		t.Error("expected 'havoc' count == 2, got", v)
	}
}
