package storage

import (
	"bytes"
	"counts/utils"
	"os"
	"path/filepath"
	"testing"
)

func initTest() {
	os.Setenv("COUNTS_DATA_DIR", "/tmp/count_data")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "data/default_config.json")
	os.Setenv("COUNTS_CONFIG", configPath)
}

func TestNoCounters(t *testing.T) {
	initTest()
	//FIXME: size of cache should be read from config
	m1 := newManager()
	m2 := newManager()
	m1.Create("marvel")
	data1 := []byte("wolverine")
	m1.SaveData("marvel", data1, 0)
	data2 := m2.LoadData("marvel", 0, 0)
	if bytes.Compare(data1, data2) != 0 {
		t.Error("Expected data2 == "+string(data1)+" got", data2)
	}
}
