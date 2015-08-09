package storage

import (
	"bufio"
	"counts/utils"
	"os"
	"path/filepath"
)

// FIXME: path currently hardcoded

func getPath(path string) string {
	var dataPath, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return dataPath
}

var dataPath = getPath("~/.counts/storage")

// the storage should deal with 2 types of on disk files, info and data
// info describes a domain and can be used to load back from disk the settings
// of a counter to reinitialize it
// the data is to refill the counters from disk
type managerStruct struct {
}

var manager *managerStruct

func getManager() *managerStruct {
	if manager == nil {
		manager = &managerStruct{}
	}
	return manager
}

/*
Create storage
*/
func (m *managerStruct) Create(ID string) {
	f, err := os.Create(filepath.Join(dataPath, ID))
	defer f.Close()
	utils.PanicOnError(err)
}

func (m *managerStruct) Save(ID string, data []byte) {
	f, err := os.Create(ID)
	defer f.Close()
	utils.PanicOnError(err)
	w := bufio.NewWriter(f)
	w.Write(data)
	w.Flush()
}

/*
Manager is responsible for manipulating the counters and syncing to disk
*/
var Manager = getManager()
