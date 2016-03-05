package config

import (
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/njpatel/loggo"

	"utils"
)

//go:generate bash bake_default_config.sh
const defaultTomlConfig = `
# This is where top level info is stored for the counter manager we
# could also use a boltDB in the DataDir but this would make it harder
# to sync over replicas since not all replicas will hold the all the
# counters.
info_dir = "~/.skizze"

# This is where the data is stored either as json or .count (pure bytes)
data_dir = "~/.skizze/data"

# The host interface for the server
host = "localhost"

# The port number for the server
port = 3596

# Treshold for saving a sketch to disk
save_threshold_seconds = 1
`

var logger = loggo.GetLogger("config")

// Config stores all configuration parameters for Go
type Config struct {
	InfoDir              string `toml:"info_dir"`
	DataDir              string `toml:"data_dir"`
	Host                 string `toml:"host"`
	Port                 int    `toml:"port"`
	SaveThresholdSeconds uint   `toml:"save_threshold_seconds"`
}

var config *Config

// MaxKeySize ...
const MaxKeySize int = 32768 // max key size BoltDB in bytes

func parseConfigTOML() *Config {
	config = &Config{}
	if _, err := toml.Decode(defaultTomlConfig, &config); err != nil {
		utils.PanicOnError(err)
	}

	configPath := os.Getenv("SKIZZE_CONFIG")
	if configPath != "" {
		_, err := os.Open(configPath)
		if err != nil {
			logger.Warningf("Unable to find config file, using defaults")
			return config
		}
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			logger.Warningf("Error parsing config file, using defaults")
		}
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

		host := strings.TrimSpace(os.Getenv("SKZ_HOST"))
		if len(host) == 0 {
			host = config.Host
		}

		port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SKZ_PORT")))
		if err != nil {
			port = config.Port
		}

		saveThresholdSecondsInt, err := strconv.Atoi(strings.TrimSpace(os.Getenv("SKZ_SAVE_TRESHOLD_SECS")))
		saveThresholdSeconds := uint(saveThresholdSecondsInt)
		if err != nil {
			saveThresholdSeconds = config.SaveThresholdSeconds
		}

		if saveThresholdSeconds < 3 {
			saveThresholdSeconds = 3
		}

		if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
			panic(err)
		}

		config = &Config{
			infoDir,
			dataDir,
			host,
			port,
			saveThresholdSeconds,
		}
	}
	return config
}

// Reset ...
func Reset() {
	config = nil
}
