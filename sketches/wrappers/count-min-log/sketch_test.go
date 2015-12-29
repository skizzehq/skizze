package cml

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
	os.RemoveAll(config.GetConfig().DataDir)
	os.RemoveAll(config.GetConfig().InfoDir)
	os.Mkdir(config.GetConfig().DataDir, 0777)
	os.Mkdir(config.GetConfig().InfoDir, 0777)
}

func TestCMLCounter(t *testing.T) {
	setupTests()
	defer tearDownTests()

	sketch, err := NewSketch(&abstract.Info{
		ID:         "avengers",
		Type:       abstract.CML,
		Properties: &abstract.Properties{},
		State:      &abstract.State{},
	})

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

	res := sketch.GetFrequency([][]byte{[]byte("cyclops")}).(map[string]uint)
	if res["cyclops"] != 3 {
		t.Error("expected 'cyclops' count == 3, got", res["cyclops"])
	}
}
