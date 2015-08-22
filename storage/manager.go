package storage

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/utils"

	"github.com/hashicorp/golang-lru"
)

var conf *config.Config
var dataPath string

// ManagerStruct the storage should deal with 2 types of on disk files, info and data
// info describes a domain and can be used to load back from disk the settings
// of a counter to reinitialize it
// the data is to refill the counters from disk
type ManagerStruct struct {
	cache *lru.Cache
}

var manager *ManagerStruct

func onFileEvicted(k interface{}, v interface{}) {
	f := v.(*os.File)
	f.Close()
}

func newManager() *ManagerStruct {
	conf = config.GetConfig()
	dataPath = conf.GetDataDir()
	cacheSize := int(conf.GetCacheSize())
	if cacheSize == 0 {
		cacheSize = 250 // default cache size
	}
	cache, err := lru.NewWithEvict(cacheSize, onFileEvicted)
	utils.PanicOnError(err)
	err = os.MkdirAll(dataPath, 0777)
	utils.PanicOnError(err)
	return &ManagerStruct{cache}
}

/*
GetManager ...
*/
func GetManager() *ManagerStruct {
	if manager == nil {
		manager = newManager()
	}
	return manager
}

/*
Create storage
*/
func (m *ManagerStruct) Create(ID string) error {
	f, err := os.Create(filepath.Join(dataPath, ID))
	if err != nil {
		return err
	}
	m.cache.Add(ID, f)
	return nil
}

/*
SaveData ...
*/
func (m *ManagerStruct) SaveData(ID string, data []byte, offset int64) error {
	f, err := m.getFileFromCache(ID)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = binary.Write(buf, binary.BigEndian, data); err != nil {
		return err
	}
	_, err = f.WriteAt(buf.Bytes(), offset)
	return err
}

/*
DeleteData ...
*/
func (m *ManagerStruct) DeleteData(ID string) error {
	v, ok := m.cache.Peek(ID)
	if ok {
		v.(*os.File).Close()
	}
	path := filepath.Join(dataPath, ID)
	/*
		if _, err := os.Stat(path); err != nil {
			return nil
		}
	*/
	return os.Remove(path)
}

/*
LoadData ...
*/
func (m *ManagerStruct) LoadData(ID string, offset int64, length int64) ([]byte, error) {
	f, err := m.getFileFromCache(ID)
	if err != nil {
		return nil, err
	}
	if length == 0 {
		if info, err := f.Stat(); err == nil {
			length = info.Size()
			length -= offset
		} else {
			return nil, err
		}
	}
	data := make([]byte, length)
	_, err = f.ReadAt(data, offset)
	return data, err
}

func (m *ManagerStruct) getFileFromCache(ID string) (*os.File, error) {
	v, ok := m.cache.Get(ID)
	var f *os.File
	var err error
	if !ok {
		if f, err = os.Open(filepath.Join(dataPath, ID)); err != nil {
			return nil, err
		}
	} else {
		f = v.(*os.File)
	}
	return f, nil
}
