package storage

import (
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

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
	infoDir := conf.InfoDir
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
