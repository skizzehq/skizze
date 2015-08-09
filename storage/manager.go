package storage

import (
	"bytes"
	"counts/utils"
	"encoding/binary"
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

func newManager() *managerStruct {
	//FIXME: size of cache should be read from config
	cache, err := lru.NewWithEvict(250, onFileEvicted)
	utils.PanicOnError(err)
	return &managerStruct{cache}
}

func getManager() *managerStruct {
	if manager == nil {
		manager = newManager()
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
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	utils.PanicOnError(err)
	_, err = f.WriteAt(buf.Bytes(), offset)
	utils.PanicOnError(err)
}

func (m *managerStruct) LoadData(ID string, offset int64, length int64) []byte {
	f := m.getFileFromCache(ID)
	if length == 0 {
		info, err := f.Stat()
		utils.PanicOnError(err)
		length = info.Size()
		length -= offset
	}
	data := make([]byte, length)
	_, err := f.ReadAt(data, offset)
	utils.PanicOnError(err)
	return data
}

func (m *managerStruct) getFileFromCache(ID string) *os.File {
	v, ok := m.cache.Get(ID)
	var f *os.File
	var err error
	if !ok {
		f, err = os.Open(filepath.Join(dataPath, ID))
		utils.PanicOnError(err)
	} else {
		f = v.(*os.File)
	}
	return f
}

func (m *managerStruct) forceFlush(ID string) {
	f := m.getFileFromCache(ID)
	m.cache.Remove(ID)
	f.Close()
}

/*
Manager is responsible for manipulating the counters and syncing to disk
*/
var Manager = getManager()
