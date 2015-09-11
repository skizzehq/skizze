package storage

import (
	"os"
	"path/filepath"
)

/*
Create storage
*/
func (m *ManagerStruct) Create(ID string) error {
	f, err := os.OpenFile(filepath.Join(dataPath, ID), os.O_CREATE|os.O_RDWR, 0644)
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
	_, err = f.WriteAt(data, offset)
	return err
}

/*
DeleteData ...
*/
func (m *ManagerStruct) DeleteData(ID string) error {
	v, ok := m.cache.Peek(ID)
	if ok {
		err := v.(*os.File).Close()
		if err != nil {
			logger.Error.Println(err)
		}
	}
	path := filepath.Join(dataPath, ID)
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
		f, err = os.OpenFile(filepath.Join(dataPath, ID), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		f = v.(*os.File)
	}
	return f, nil
}
