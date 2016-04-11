package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"utils"
	"config"
)

// SetupTests ...
func SetupTests() {
	dataDir := fmt.Sprintf("/tmp/skizze_storage_data_test")
	infoDir := fmt.Sprintf("/tmp/skizze_storage_info_test")

	// Cleanup any previously aborted test runs
	utils.PanicOnError(os.RemoveAll(dataDir))
	utils.PanicOnError(os.RemoveAll(infoDir))

	utils.PanicOnError(os.Setenv("SKIZZE_DATA_DIR", dataDir))
	utils.PanicOnError(os.Setenv("SKIZZE_INFO_DIR", infoDir))

	utils.PanicOnError(os.Mkdir(os.Getenv("SKIZZE_DATA_DIR"), 0777))
	utils.PanicOnError(os.Mkdir(os.Getenv("SKIZZE_INFO_DIR"), 0777))

	path, err := os.Getwd()
	utils.PanicOnError(err)
	path = filepath.Dir(path)
	configPath := filepath.Join(path, "config/default.toml")
	utils.PanicOnError(os.Setenv("SKIZZE_CONFIG", configPath))
	config.DataDir = dataDir
	config.InfoDir = infoDir
}

// TearDownTests ...
func TearDownTests() {
	utils.PanicOnError(os.RemoveAll(os.Getenv("SKIZZE_DATA_DIR")))
	utils.PanicOnError(os.RemoveAll(os.Getenv("SKIZZE_INFO_DIR")))
	time.Sleep(50 * time.Millisecond)
	config.Reset()
}
