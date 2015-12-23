package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/seiflotfy/skizze/utils"

	"github.com/BurntSushi/toml"
)

// Config stores all configuration parameters for Go
type Config struct {
	InfoDir              string `toml:"info_dir"`
	DataDir              string `toml:"data_dir"`
	SliceSize            uint   `toml:"slice_size"`
	CacheSize            uint   `toml:"cache_size"`
	SliceCacheSize       uint   `toml:"slice_cache_size"`
	Port                 uint   `toml:"port"`
	SaveThresholdSeconds uint   `toml:"save_threshold_seconds"`
	SaveThresholdOps     uint   `toml:"save_threshold_ops"`
}

var config *Config

// MaxKeySize ...
const MaxKeySize int = 32768 // max key size BoltDB in bytes

func parseConfigTOML() *Config {
	configPath := os.Getenv("SKZ_CONFIG")
	if configPath == "" {
		path, err := os.Getwd()
		utils.PanicOnError(err)
		path, err = filepath.Abs(path)
		utils.PanicOnError(err)
		configPath = filepath.Join(path, "config/default.toml")
	}
	_, err := os.Open(configPath)
	utils.PanicOnError(err)
	config = &Config{}
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		utils.PanicOnError(err)
	}
	return config
}

// GetConfig returns a singleton Configuration
func GetConfig() *Config {
	if config == nil {
		config = parseConfigTOML()
		usr, err := user.Current()
		utils.PanicOnError(err)
		dir := usr.HomeDir

		infoDir := strings.TrimSpace(os.Getenv("SKZ_INFO_DIR"))
		if len(infoDir) == 0 {
			if config.InfoDir[:2] == "~/" {
				infoDir = strings.Replace(config.InfoDir, "~", dir, 1)
			}
		}

		dataDir := strings.TrimSpace(os.Getenv("SKZ_DATA_DIR"))
		if len(dataDir) == 0 {
			if config.DataDir[:2] == "~/" {
				dataDir = strings.Replace(config.DataDir, "~", dir, 1)
			}
		}

		portInt, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SKZ_PORT")))
		port := uint(portInt)
		if err != nil {
			port = config.Port
		}

		saveThresholdSecondsInt, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SKZ_SAVE_TRESHOLD_SECS")))
		saveThresholdSeconds := uint(saveThresholdSecondsInt)
		if err != nil {
			saveThresholdSeconds = config.SaveThresholdSeconds
		}

		saveThresholdOpsInt, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SKZ_SAVE_TRESHOLD_OPS")))
		saveThresholdOps := uint(saveThresholdOpsInt)
		if err != nil {
			saveThresholdOps = config.SaveThresholdOps
		}

		config = &Config{
			infoDir,
			dataDir,
			config.SliceSize,
			config.CacheSize,
			config.SliceCacheSize,
			port,
			saveThresholdSeconds,
			saveThresholdOps,
		}
	}
	return config
}
