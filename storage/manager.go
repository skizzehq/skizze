package storage

import (
	"os"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/utils"

	"github.com/hashicorp/golang-lru"
)

var conf *config.Config
var dataPath string
var logger = utils.GetLogger()

// ManagerStruct the storage should deal with 2 types of on disk files, info and data
// info describes a domain and can be used to load back from disk the settings
// of a counter to reinitialize it
// the data is to refill the counters from disk
type ManagerStruct struct {
	cache *lru.Cache
}

var manager *ManagerStruct

func newManager() *ManagerStruct {
	conf = config.GetConfig()
	dataPath = conf.DataDir
	cacheSize := int(conf.CacheSize)
	if cacheSize == 0 {
		cacheSize = 250 // default cache size
	}
	cache, err := lru.NewWithEvict(cacheSize, func(k interface{}, v interface{}) {
		f := v.(*os.File)
		err := f.Close()
		if err != nil {
			logger.Error.Println(err)
		}
	})
	utils.PanicOnError(err)
	err = os.MkdirAll(dataPath, 0777)
	utils.PanicOnError(err)
	return &ManagerStruct{cache}
}

/*
Manager ...
*/
func Manager() *ManagerStruct {
	if manager == nil {
		manager = newManager()
	}
	return manager
}
