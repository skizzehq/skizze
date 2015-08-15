package cuckoofilter

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/seiflotfy/counts/config"
	"github.com/seiflotfy/counts/utils"
)

func setupTests() {
	os.Setenv("COUNTS_DATA_DIR", "/tmp/count_data")
	os.Setenv("COUNTS_INFO_DIR", "/tmp/count_info")
	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "../../../config/default.toml")
	os.Setenv("COUNTS_CONFIG", configPath)
	tearDownTests()
}

func tearDownTests() {
	os.RemoveAll(config.GetConfig().GetDataDir())
	os.RemoveAll(config.GetConfig().GetInfoDir())
	os.Mkdir(config.GetConfig().GetDataDir(), 0777)
	os.Mkdir(config.GetConfig().GetInfoDir(), 0777)
}

func TestInsertion(t *testing.T) {
	setupTests()
	//defer tearDownTests()

	cf := NewCuckooFilter("ultimates", 1000000)

	fd, err := os.Open("/usr/share/dict/web2")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	var values [][]byte
	for scanner.Scan() {
		s := []byte(scanner.Text())
		cf.InsertUnique(s)
		values = append(values, s)
	}

	count := cf.GetCount()
	if count != 235081 {
		t.Errorf("Expected count = 235081, instead count = %d", count)
	}

	for _, v := range values {
		cf.Delete(v)
	}

	count = cf.GetCount()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}
