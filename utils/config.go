package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type configStruct struct {
	infoDir        string
	dataDir        string
	sliceSize      uint
	cacheSize      uint
	sliceCacheSize uint
}

var config *configStruct

/*
GetInfoDir returns the top level info
*/
func (c *configStruct) GetInfoDir() string {
	return c.infoDir
}

/*
GetDataDir returns the top level info
*/
func (c *configStruct) GetDataDir() string {
	return c.dataDir
}

/*
GetSliceSize returns the top level info
*/
func (c *configStruct) GetSliceSize() uint {
	return c.sliceSize
}

/*
GetCacheSize returns the top level info
*/
func (c *configStruct) GetCacheSize() uint {
	return c.cacheSize
}

/*
GetSliceCacheSize returns the top level info
*/
func (c *configStruct) GetSliceCacheSize() uint {
	return c.sliceCacheSize
}

/*
GetConfig returns a singleton Configuration
*/
func getConfig() *configStruct {
	if config == nil {
		configPath := os.Getenv("COUNTS_CONFIG")
		if configPath == "" {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				logger.Error.Println("error:", err)
			}
			configPath = filepath.Join(dir, "data/default_config.json")
		}
		file, _ := os.Open(configPath)
		decoder := json.NewDecoder(file)
		tempConfig := &tempConfigStruct{}
		err := decoder.Decode(&tempConfig)
		config = &configStruct{
			tempConfig.InfoDir,
			tempConfig.DataDir,
			tempConfig.SliceSize,
			tempConfig.CacheSize,
			tempConfig.SliceCacheSize,
		}
		fmt.Println(config)
		if err != nil {
			logger.Error.Println("error:", err)
		}
	}

	return config
}

/*
Config
*/
var Config = getConfig()
