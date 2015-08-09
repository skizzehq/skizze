package storage

import (
	"counts/utils"
	"os"
	"os/user"
	"path/filepath"

	"github.com/hashicorp/golang-lru"
)

//FIXME: path currently hardcoded

func getPath(path string) string {
	const storeDir = ".counts/data"
	usr, _ := user.Current()
	dataPath := filepath.Join(usr.HomeDir, storeDir)
	return dataPath
}

var dataPath = getPath("~/.counts/storage")

// the storage should deal with 2 types of on disk files, info and data
// info describes a domain and can be used to load back from disk the settings
// of a counter to reinitialize it
// the data is to refill the counters from disk
type managerStruct struct {
	cache *lru.Cache
}

var manager *managerStruct

func onFileEvicted(k interface{}, v interface{}) {
	f := v.(*os.File)
	f.Close()
}

func getManager() *managerStruct {
	if manager == nil {
		//FIXME: size of cache should be read from config
		cache, err := lru.NewWithEvict(250, onFileEvicted)
		utils.PanicOnError(err)
		manager = &managerStruct{cache}
	}
	return manager
}

/*
Create storage
*/
func (m *managerStruct) Create(ID string) {
	f, err := os.Create(filepath.Join(dataPath, ID))
	utils.PanicOnError(err)
	m.cache.Add(ID, f)
}

func (m *managerStruct) SaveData(ID string, data []byte, offset int64) {
	f := m.getFileFromCache(ID)
	_, err := f.WriteAt(data, offset)
	utils.PanicOnError(err)
}

func (m *managerStruct) LoadData(ID string, offset int64) []byte {
	f := m.getFileFromCache(ID)
	data := []byte{}
	_, err := f.ReadAt(data, offset)
	utils.PanicOnError(err)
	return data
}

func (m *managerStruct) getFileFromCache(ID string) *os.File {
	v, ok := m.cache.Get(ID)
	var f *os.File
	var err error
	if !ok {
		f, err = os.Create(ID)
		utils.PanicOnError(err)
	} else {
		f = v.(*os.File)
	}
	return f
}

/*
Manager is responsible for manipulating the counters and syncing to disk
*/
var Manager = getManager()
