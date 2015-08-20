package storage

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"time"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/utils"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/golang-lru"
)

var conf *config.Config
var dataPath string
var db *bolt.DB

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

/*
LoadAllInfo ...
*/
func (m *ManagerStruct) LoadAllInfo() ([][]byte, error) {
	db, err := getInfoDB()
	if err != nil {
		return nil, err
	}
	var infoDatas [][]byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("info"))

		infoDatas = make([][]byte, bucket.Stats().KeyN)
		c := bucket.Cursor()
		i := 0
		for k, infoData := c.First(); k != nil; k, infoData = c.Next() {
			// We need to copy as the infoData is freed once
			// we are done with this transaction
			infoDataResult := make([]byte, len(infoData))
			copy(infoDataResult, infoData)
			infoDatas[i] = infoDataResult
			i++
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	return infoDatas, nil
}

/*
SaveInfo ...
*/
func (m *ManagerStruct) SaveInfo(id string, infoData []byte) error {
	db, err := getInfoDB()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("info"))
		if err != nil {
			return err
		}
		key := []byte(id)
		err = bucket.Put(key, infoData)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

/*
DeleteInfo ...
*/
func (m *ManagerStruct) DeleteInfo(id string) error {
	db, err := getInfoDB()
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("info"))
		if err != nil {
			return err
		}
		key := []byte(id)
		err = bucket.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

/*
CloseInfoDB ...
*/
func CloseInfoDB() error {
	if db != nil {
		err := db.Close()
		db = nil
		return err
	}
	return nil
}

/*
getInfoDB returns a singleton of the infoDB
*/
func getInfoDB() (*bolt.DB, error) {
	if db != nil {
		return db, nil
	}
	var err error
	infoDir := conf.GetInfoDir()
	err = os.MkdirAll(infoDir, 0777)
	if err != nil {
		return nil, err
	}
	infoDBPath := filepath.Join(infoDir, "info.db")
	dbOptions := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(infoDBPath, 0600, dbOptions)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("info"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
