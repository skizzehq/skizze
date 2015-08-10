package utils

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

/*
Configuration stores all configuration parameters for Go
*/
type tempConfigStruct struct {
	InfoDir        string `json:"InfoDir"`
	DataDir        string `json:"DataDir"`
	SliceSize      uint   `json:"SliceSize"`
	CacheSize      uint   `json:"CacheSize"`
	SliceCacheSize uint   `json:"SliceCacheSize"`
	Port           int    `json:"Port"`
}

/*
ConfigStruct ...
*/
type ConfigStruct struct {
	infoDir        string
	dataDir        string
	sliceSize      uint
	cacheSize      uint
	sliceCacheSize uint
	port           int
}

var config *ConfigStruct

/*
GetInfoDir returns the top level info
*/
func (c *ConfigStruct) GetInfoDir() string {
	return c.infoDir
}

/*
GetDataDir returns the top level info
*/
func (c *ConfigStruct) GetDataDir() string {
	return c.dataDir
}

/*
GetSliceSize returns the top level info
*/
func (c *ConfigStruct) GetSliceSize() uint {
	return c.sliceSize
}

/*
GetCacheSize returns the top level info
*/
func (c *ConfigStruct) GetCacheSize() uint {
	return c.cacheSize
}

/*
GetSliceCacheSize returns the top level info
*/
func (c *ConfigStruct) GetSliceCacheSize() uint {
	return c.sliceCacheSize
}

/*
GetPort returns the port the server runs on
*/
func (c *ConfigStruct) GetPort() int {
	return c.port
}

func parseConfigJSON() *tempConfigStruct {
	configPath := os.Getenv("COUNTS_CONFIG")
	if configPath == "" {
		path, err := os.Getwd()
		PanicOnError(err)
		path, err = filepath.Abs(path)
		PanicOnError(err)
		configPath = filepath.Join(path, "data/default_config.json")
	}
	file, err := os.Open(configPath)
	PanicOnError(err)
	decoder := json.NewDecoder(file)
	tempConfig := &tempConfigStruct{}
	err = decoder.Decode(&tempConfig)
	PanicOnError(err)
	return tempConfig
}

/*
GetConfig returns a singleton Configuration
*/
func GetConfig() *ConfigStruct {
	if config == nil {
		tempConfig := parseConfigJSON()
		usr, err := user.Current()
		PanicOnError(err)
		dir := usr.HomeDir

		infoDir := strings.TrimSpace(os.Getenv("COUNTS_INFO_DIR"))
		if len(infoDir) == 0 {
			if tempConfig.InfoDir[:2] == "~/" {
				infoDir = strings.Replace(tempConfig.InfoDir, "~", dir, 1)
			}
		}
		os.Mkdir(infoDir, 0777)

		dataDir := strings.TrimSpace(os.Getenv("COUNTS_DATA_DIR"))
		if len(dataDir) == 0 {
			if tempConfig.DataDir[:2] == "~/" {
				dataDir = strings.Replace(tempConfig.DataDir, "~", dir, 1)
			}
		}
		os.Mkdir(dataDir, 0777)

		port, err := strconv.Atoi(strings.TrimSpace(os.Getenv("COUNTS_PORT")))
		if err != nil {
			port = tempConfig.Port
		}

		config = &ConfigStruct{
			infoDir,
			dataDir,
			tempConfig.SliceSize,
			tempConfig.CacheSize,
			tempConfig.SliceCacheSize,
			port,
		}
	}
	return config
}
