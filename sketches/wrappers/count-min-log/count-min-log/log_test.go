package cml

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/utils"
)

func setupTests() {
	os.Setenv("SKZ_DATA_DIR", "/tmp/skizze_data")
	os.Setenv("SKZ_INFO_DIR", "/tmp/skizze_info")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "../../../config/default.toml")
	os.Setenv("SKZ_CONFIG", configPath)
	tearDownTests()
}

func tearDownTests() {
	os.RemoveAll(config.GetConfig().GetDataDir())
	os.RemoveAll(config.GetConfig().GetInfoDir())
	os.Mkdir(config.GetConfig().GetDataDir(), 0777)
	os.Mkdir(config.GetConfig().GetInfoDir(), 0777)
}

// Ensures that Add adds to the set and Count returns the correct
// approximation.
func TestLog16AddAndCount(t *testing.T) {
	setupTests()
	defer tearDownTests()
	info := &abstract.Info{ID: "ultimates",
		Type:       abstract.CML,
		Properties: map[string]float64{},
		State:      make(map[string]uint64)}

	log, _ := NewForCapacity16(info, 1000, 0.01)

	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))
	if count := log.GetCount([]byte("a")); uint(count) != 3 {
		t.Errorf("expected 3, got %d", uint(count))
	}

	if count := log.GetCount([]byte("b")); uint(count) != 2 {
		t.Errorf("expected 2, got %d", uint(count))
	}

	if count := log.GetCount([]byte("c")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := log.GetCount([]byte("d")); uint(count) != 1 {
		t.Errorf("expected 1, got %d", uint(count))
	}

	if count := log.GetCount([]byte("x")); uint(count) != 0 {
		t.Errorf("expected 0, got %d", uint(count))
	}
}

// Ensures that Reset restores the sketch to its original state.
func TestLog16Reset(t *testing.T) {
	setupTests()
	defer tearDownTests()
	info := &abstract.Info{ID: "ultimates",
		Type:       abstract.CML,
		Properties: map[string]float64{},
		State:      make(map[string]uint64)}
	log, _ := NewForCapacity16(info, 1000, 0.001)
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("c"))
	log.IncreaseCount([]byte("b"))
	log.IncreaseCount([]byte("d"))
	log.IncreaseCount([]byte("a"))
	log.IncreaseCount([]byte("a"))

	log.Reset()

	for i := uint(0); i < log.k; i++ {
		for j := uint(0); j < log.w; j++ {
			if x := log.store[i][j]; x != 0 {
				t.Errorf("expected matrix to be completely empty, got %d", x)
			}
		}
	}
}

/*
func TestLog16Reset(t *testing.T) {
	setupTests()
	defer tearDownTests()
	info := &abstract.Info{ID: "ultimates",
		Type:       abstract.CML,
		Properties: map[string]float64{},
		State:      make(map[string]uint64)}
	log, _ := NewForCapacity16(info, 10, 0.9)
	log.store[0][0] = 1
	log.store[0][2] = 1
	log.store[1][2] = 1
	//fmt.Println(log.store)
	//fmt.Println("======") //, len(log.store), len(log.store[0]))
	log.registers.save(log.store)
	//fmt.Println("=====================")
	log.store, _ = log.registers.load()
}
*/
