package storage

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
)

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
