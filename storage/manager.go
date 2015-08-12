package storage

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/seiflotfy/counts/config"
	"github.com/seiflotfy/counts/utils"

	"github.com/hashicorp/golang-lru"
)

//FIXME: path currently hardcoded

func getPath(path string) string {
	const storeDir = ".counts/data"
	usr, _ := user.Current()
	dataPath := filepath.Join(usr.HomeDir, storeDir)
	return dataPath
}

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
	//FIXME: size of cache should be read from config
	cache, err := lru.NewWithEvict(250, onFileEvicted)
	os.MkdirAll(dataPath, 0777)
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
func (m *ManagerStruct) Create(ID string) {
	f, err := os.Create(filepath.Join(dataPath, ID))
	utils.PanicOnError(err)
	m.cache.Add(ID, f)
}

/*
SaveData ...
*/
func (m *ManagerStruct) SaveData(ID string, data []byte, offset int64) {
	f := m.getFileFromCache(ID)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, data)
	utils.PanicOnError(err)
	_, err = f.WriteAt(buf.Bytes(), offset)
	utils.PanicOnError(err)
}

/*
LoadData ...
*/
func (m *ManagerStruct) LoadData(ID string, offset int64, length int64) []byte {
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

func (m *ManagerStruct) getFileFromCache(ID string) *os.File {
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

func (m *ManagerStruct) forceFlush(ID string) {
	f := m.getFileFromCache(ID)
	m.cache.Remove(ID)
	f.Close()
}

/*
LoadAllInfo ...
*/
func (m *ManagerStruct) LoadAllInfo() [][]byte {
	infoDir := conf.GetInfoDir()
	err := os.MkdirAll(infoDir, 0777)
	fileInfos, err := ioutil.ReadDir(infoDir)
	utils.PanicOnError(err)
	infoDatas := make([][]byte, len(fileInfos))
	for i, fileInfo := range fileInfos {
		filePath := filepath.Join(infoDir, fileInfo.Name())
		infoData, err := ioutil.ReadFile(filePath)
		utils.PanicOnError(err)

		infoDatas[i] = infoData
	}

	return infoDatas
}

/*
SaveInfo ...
*/
func (m *ManagerStruct) SaveInfo(id string, infoData []byte) {
	infoDir := conf.GetInfoDir()
	infoPath := filepath.Join(infoDir, id+".json")
	f, err := os.Create(infoPath)
	defer f.Close()
	utils.PanicOnError(err)
	f.Write(infoData)
}
