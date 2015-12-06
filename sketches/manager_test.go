package sketches

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/sketches/wrappers/topk"
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
	os.Setenv("SKZ_SAVE_TRESHOLD_OPS", "1")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	os.Setenv("SKZ_CONFIG", configPath)
	tearDownTests()
}

func tearDownTests() {
	storage.CloseInfoDB()
	os.RemoveAll(config.GetConfig().DataDir)
	os.RemoveAll(config.GetConfig().InfoDir)
	os.Mkdir(config.GetConfig().DataDir, 0777)
	os.Mkdir(config.GetConfig().InfoDir, 0777)
	manager.Destroy()
}

func TestNoSketches(t *testing.T) {
	setupTests()
	defer tearDownTests()
	var manager, err = newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	var sketches []string
	sketches, err = manager.GetSketches()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
	}
}

func TestDuplicateSketches(t *testing.T) {
	setupTests()
	defer tearDownTests()
	var manager, err = newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	props := map[string]float64{"capacity": 10000000.0}

	err = manager.CreateSketch("marvel", "hllpp", props)
	if err != nil {
		t.Error("Expected no errors while creating sketch, got", err)
	}
	err = manager.CreateSketch("marvel", "hllpp", props)
	if err == nil {
		t.Error("Expected errors while creating sketch duplicate sketch, got", err)
	}
}

func TestDefaultCounter(t *testing.T) {
	setupTests()
	defer tearDownTests()

	var manager, err = newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	var sketches []string
	sketches, err = manager.GetSketches()
	if err != nil {
		t.Error("Expected no errors while getting sketches, got", err)
	}
	if len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
	}

	props := map[string]float64{"capacity": 10000000.0}
	err = manager.CreateSketch("marvel", "hllpp", props)
	if err != nil {
		t.Error("Expected no errors while creating sketch, got", err)
	}

	sketches, err = manager.GetSketches()
	if err != nil {
		t.Error("Expected no errors while getting sketches, got", err)
	}
	if len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	}

	err = manager.AddToSketch("marvel", "hllpp", []string{"hulk", "thor"})
	if err != nil {
		t.Error("Expected no errors while adding to sketch, got", err)
	}

	var res map[string]interface{}
	res, err = manager.GetCountForSketch("marvel", "hllpp", nil)
	if len(sketches) != 1 {
		t.Error("Expected 1 sketches, got", len(sketches))
	}

	if res["result"].(uint) != 2 {
		t.Error("Expected count == 2, got", res["count"].(uint))
	}

	err = manager.DeleteSketch("marvel", "hllpp")
	if err != nil {
		t.Error("Expected no errors while deleting sketch, got", err)
	}

	sketches, err = manager.GetSketches()
	if err != nil {
		t.Error("Expected no errors while getting sketches, got", err)
	}
	if len(sketches) != 0 {
		t.Error("Expected 0 sketches, got", len(sketches))
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
	if _, exists = m1.info["x-force.hllpp"]; exists {
		t.Error("expected x-force to not be initially loaded by manager")
	}

	props := map[string]float64{"capacity": 10000000.0}
	err = m1.CreateSketch("x-force", "hllpp", props)
	if err != nil {
		t.Fatal(err)
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	if _, exists = m2.info["x-force.hllpp"]; !exists {
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

	props := map[string]float64{"capacity": 10000000.0}
	m1.CreateSketch("avengers", "hllpp", props)

	m1.AddToSketch("avengers", "hllpp", []string{"sabertooth",
		"thunderbolt", "havoc", "cyclops"})

	res, err := m1.GetCountForSketch("avengers", "hllpp", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	if res["result"].(uint) != 4 {
		t.Error("expected avengers to have count 4, got", res["result"].(uint))
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	res, err = m2.GetCountForSketch("avengers", "hllpp", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	if res["result"].(uint) != 4 {
		t.Error("expected avengers to have count 4, got", res["result"].(uint))
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

	props := map[string]float64{"capacity": 10000000.0}
	m1.CreateSketch("avengers", "hllpp", props)
	m1.CreateSketch("x-men", "hllpp", props)

	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		t.Error(err)
	}
	scanner := bufio.NewScanner(fd)

	i := 0
	values := []string{} //{"a", "aam"} //"doorknob", "doorless"
	for scanner.Scan() {
		s := []byte(scanner.Text())
		values = append(values, string(s))
		i++
		if i == 1000 {
			break
		}
	}

	// Add all values in a go routine per value
	var wg sync.WaitGroup
	defer wg.Wait()
	resChan := make(chan interface{})

	var pFunc = func(value string) {
		defer wg.Done()
		m1.AddToSketch("avengers", "hllpp", []string{value})
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
	m1.AddToSketch("x-men", "hllpp", values)
	count1, err := m1.GetCountForSketch("avengers", "hllpp", nil)
	count2, err := m1.GetCountForSketch("x-men", "hllpp", nil)
	if count1["result"].(uint) != count2["result"].(uint) {
		t.Error("expected avengers count == x-men count, got", count1["result"].(uint), "!=", count2["result"].(uint))
	}
}

func TestFailCreateSketch(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}

	// test for unknown sketchType
	props := map[string]float64{"capacity": 10000000.0}
	err = m1.CreateSketch("marvel", "wrong", props)
	if err == nil {
		t.Error("Expected errors while creating sketch, got", err)
	}

	buffer := make([]byte, config.MaxKeySize+1)
	for i := 0; i < config.MaxKeySize+1; i++ {
		buffer[i] = byte(49) // ascii 1
	}
	sketchID := string(buffer)
	// test for too long sketchID
	err = m1.CreateSketch(sketchID, "hllpp", props)
	if err == nil {
		t.Error("Expected errors while creating sketch, got", err)
	}
}

func TestFailDeleteSketch(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}
	err = m1.DeleteSketch("-1", "hllpp")
	if err == nil {
		t.Error("Expected error, got", err)
	}
}

func TestFailDeleteFromSketch(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}
	err = m1.DeleteFromSketch("-1", "hlpp", []string{})
	if err == nil {
		t.Error("Expected error, got", err)
	}
}

func TestFailGetCountForSketch(t *testing.T) {
	setupTests()
	defer tearDownTests()

	m1, err := newManager()
	if err != nil {
		t.Log("Expected no errors, got", err)
	}
	_, err = m1.GetCountForSketch("-1", "hllpp", nil)
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
	props := map[string]float64{"capacity": 3.0}
	m1.CreateSketch("avengers", "topk", props)

	err = m1.AddToSketch("avengers", "topk", []string{"sabertooth",
		"thunderbolt", "havoc", "cyclops", "cyclops", "cyclops", "havoc"})

	res, err := m1.GetCountForSketch("avengers", "topk", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	top := res["result"].([]topk.ResultElement)
	if len(top) != 3 {
		t.Error("expected avengers to have 3 elements, got", len(top))
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}
	res, err = m2.GetCountForSketch("avengers", "topk", nil)
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	top = res["result"].([]topk.ResultElement)
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
	props := map[string]float64{"epsilon": 0.5}
	m1.CreateSketch("avengers", abstract.CML, props)

	m1.AddToSketch("avengers", "cml", []string{"sabertooth",
		"thunderbolt", "havoc", "cyclops", "cyclops", "cyclops", "havoc"})

	_, err = m1.GetCountForSketch("avengers", "cml", []string{"cyclops"})
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	m2, err := newManager()
	if err != nil {
		t.Error("Expected no errors, got", err)
	}

	var res map[string]interface{}
	res, err = m2.GetCountForSketch("avengers", "cml", []string{"cyclops"})
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	counts := res["result"].(map[string]uint)
	if v, ok := counts["cyclops"]; !ok {
		t.Error("expected to find 'cyclops' in avengers, got", ok)
	} else if v != 3 {
		t.Error("expected 'cyclops' count == 3, got", v)
	}

	res, err = m2.GetCountForSketch("avengers", "cml", []string{"havoc"})
	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}
	counts = res["result"].(map[string]uint)
	if v, ok := counts["havoc"]; !ok {
		t.Error("expected to find 'havoc' in avengers, got", ok)
	} else if v != 2 {
		t.Error("expected 'havoc' count == 2, got", v)
	}
}
