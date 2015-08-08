package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

/*
Configuration stores all configuration parameters for Go
*/
type configurationStruct struct {
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

var config *configurationStruct
var logg = GetLogger()

/*
GetInfoDir returns the top level info
*/
func (c configurationStruct) GetInfoDir() string {
	return config.InfoDir
}

/*
GetDataDir returns the top level info
*/
func (c configurationStruct) GetDataDir() string {
	return config.DataDir
}

/*
GetSliceSize returns the top level info
*/
func (c configurationStruct) GetSliceSize() uint {
	return config.SliceSize
}

/*
GetCacheSize returns the top level info
*/
func (c configurationStruct) GetCacheSize() uint {
	return config.CacheSize
}

/*
GetSliceCacheSize returns the top level info
*/
func (c configurationStruct) GetSliceCacheSize() uint {
	return config.SliceCacheSize
}

/*
GetConfig returns a singleton Configuration
*/
func getConfig() *configurationStruct {
	if config == nil {
		configPath := os.Getenv("COUNTS_CONFIG")
		if configPath == "" {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				logg.Error.Println("error:", err)
			}
			configPath = dir + "/data/default_config.json"
		}
		file, _ := os.Open(configPath)
		decoder := json.NewDecoder(file)
		config = &configurationStruct{}
		err := decoder.Decode(&config)
		if err != nil {
			logg.Error.Println("error:", err)
		}
	}

	return config
}

/*
Config
*/
var Config = getConfig()
