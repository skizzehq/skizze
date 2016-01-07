package utils

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// SetupTests ...
func SetupTests() {
	rand.Seed(time.Now().Unix())
	dataDir := fmt.Sprintf("/tmp/skizze_storage_data_%d", rand.Uint32())
	infoDir := fmt.Sprintf("/tmp/skizze_storage_info_%d", rand.Uint32())
	PanicOnError(os.Setenv("SKZ_DATA_DIR", dataDir))
	PanicOnError(os.Setenv("SKZ_INFO_DIR", infoDir))

	PanicOnError(os.Mkdir(os.Getenv("SKZ_DATA_DIR"), 0777))
	PanicOnError(os.Mkdir(os.Getenv("SKZ_INFO_DIR"), 0777))

	path, err := os.Getwd()
	PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	PanicOnError(os.Setenv("SKZ_CONFIG", configPath))
}

// TearDownTests ...
func TearDownTests() {
	PanicOnError(os.RemoveAll(os.Getenv("SKZ_DATA_DIR")))
	PanicOnError(os.RemoveAll(os.Getenv("SKZ_INFO_DIR")))
}
