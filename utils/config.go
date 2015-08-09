package utils

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

/*
Configuration stores all configuration parameters for Go
*/
type tempConfigStruct struct {
	// this is where top level info is stored for the counter manager we could also use a boltDB in the DataDir but this would make it harder to sync over replicas since not all replicas will hold the all the counters.
	InfoDir string `json:"InfoDir"`
	// this is where the data is stored either as json or .count (pure bytes)
	DataDir string `json:"DataDir"`
	// is needed for the raw bytes storage since we can split them up and not have it all in memory at once.
	SliceSize uint `json:"SliceSize"`
	// num of counters in cache
	CacheSize uint `json:"CacheSize"`
	// number slices in the slice cache
	SliceCacheSize uint `json:"SliceCacheSize"`
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
GetConfig returns a singleton Configuration
*/
func GetConfig() *ConfigStruct {
	if config == nil {
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

		config = &ConfigStruct{
			infoDir,
			dataDir,
			tempConfig.SliceSize,
			tempConfig.CacheSize,
			tempConfig.SliceCacheSize,
		}
		if err != nil {
			logger.Error.Println("error:", err)
		}
	}
	return config
}
