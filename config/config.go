package config

import (
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/seiflotfy/counts/utils"

	"github.com/BurntSushi/toml"
)

/*
Config stores all configuration parameters for Go
*/
type Config struct {
	InfoDir        string `toml:"info_dir"`
	DataDir        string `toml:"data_dir"`
	SliceSize      uint   `toml:"slice_size"`
	CacheSize      uint   `toml:"cache_size"`
	SliceCacheSize uint   `toml:"slice_cache_size"`
	Port           int    `toml:"port"`
}

var config *Config

/*
GetInfoDir returns the top level info
*/
func (c *Config) GetInfoDir() string {
	return c.InfoDir
}

/*
GetDataDir returns the top level info
*/
func (c *Config) GetDataDir() string {
	return c.DataDir
}

/*
GetSliceSize returns the top level info
*/
func (c *Config) GetSliceSize() uint {
	return c.SliceSize
}

/*
GetCacheSize returns the top level info
*/
func (c *Config) GetCacheSize() uint {
	return c.CacheSize
}

/*
GetSliceCacheSize returns the top level info
*/
func (c *Config) GetSliceCacheSize() uint {
	return c.SliceCacheSize
}

/*
GetPort returns the port the server runs on
*/
func (c *Config) GetPort() int {
	return c.Port
}

func parseConfigTOML() *Config {
	configPath := os.Getenv("COUNTS_CONFIG")
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

/*
GetConfig returns a singleton Configuration
*/
func GetConfig() *Config {
	if config == nil {
		config = parseConfigTOML()
		usr, err := user.Current()
		utils.PanicOnError(err)
		dir := usr.HomeDir

		infoDir := strings.TrimSpace(os.Getenv("COUNTS_INFO_DIR"))
		if len(infoDir) == 0 {
			if config.InfoDir[:2] == "~/" {
				infoDir = strings.Replace(config.InfoDir, "~", dir, 1)
			}
		}
		err = os.Mkdir(infoDir, 0777)
		utils.PanicOnError(err)

		dataDir := strings.TrimSpace(os.Getenv("COUNTS_DATA_DIR"))
		if len(dataDir) == 0 {
			if config.DataDir[:2] == "~/" {
				dataDir = strings.Replace(config.DataDir, "~", dir, 1)
			}
		}
		err = os.Mkdir(dataDir, 0777)
		utils.PanicOnError(err)

		port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("COUNTS_PORT")))
		if err != nil {
			port = config.Port
		}

		config = &Config{
			infoDir,
			dataDir,
			config.SliceSize,
			config.CacheSize,
			config.SliceCacheSize,
			port,
		}
	}
	return config
}
