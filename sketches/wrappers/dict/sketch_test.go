package dict

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/sketches/abstract"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

func setupTests() {
	os.Setenv("SKZ_DATA_DIR", "/tmp/skizze_manager_data")
	os.Setenv("SKZ_INFO_DIR", "/tmp/skizze_manager_info")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "../../config/default.toml")
	os.Setenv("SKZ_CONFIG", configPath)
	tearDownTests()
}

func tearDownTests() {
	storage.CloseInfoDB()
	os.RemoveAll(config.GetConfig().GetDataDir())
	os.RemoveAll(config.GetConfig().GetInfoDir())
	os.Mkdir(config.GetConfig().GetDataDir(), 0777)
	os.Mkdir(config.GetConfig().GetInfoDir(), 0777)
}

func TestAddMultiple(t *testing.T) {
	setupTests()
	defer tearDownTests()

	sketch, err := NewSketch(&abstract.Info{
		ID:         "avengers",
		Type:       abstract.Dict,
		Properties: make(map[string]float64),
		State:      make(map[string]uint64)})

	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	values := [][]byte{
		[]byte("sabertooth"),
		[]byte("thunderbolt"),
		[]byte("havoc"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("havoc")}

	sketch.AddMultiple(values)

	res := sketch.GetFrequency([][]byte{[]byte("cyclops")}).(map[string]int)
	if res["cyclops"] != 3 {
		t.Error("expected 'cyclops' count == 3, got", res["cyclops"])
	}
}

func TestDecrease(t *testing.T) {
	setupTests()
	defer tearDownTests()

	sketch, err := NewSketch(&abstract.Info{
		ID:         "avengers",
		Type:       abstract.Dict,
		Properties: make(map[string]float64),
		State:      make(map[string]uint64)})

	if err != nil {
		t.Error("expected avengers to have no error, got", err)
	}

	addValues := [][]byte{
		[]byte("sabertooth"),
		[]byte("thunderbolt"),
		[]byte("havoc"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("cyclops"),
		[]byte("havoc")}

	sketch.AddMultiple(addValues)

	removeValues := [][]byte{
		[]byte("sabertooth"),
		[]byte("thunderbolt"),
		[]byte("havoc"),
		[]byte("cyclops"),
		[]byte("havoc")}

	sketch.RemoveMultiple(removeValues)

	res := sketch.GetFrequency([][]byte{[]byte("cyclops")}).(map[string]int)
	if res["cyclops"] != 2 {
		t.Error("expected 'cyclops' count == 2, got", res["cyclops"])
	}
}
