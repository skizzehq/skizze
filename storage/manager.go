package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skizzehq/skizze/config"
	"github.com/skizzehq/skizze/datamodel"
	"github.com/skizzehq/skizze/utils"
)

// Manager the storage should deal with 2 types of on disk files, info and data
// info describes a domain and can be used to load back from disk the settings
// of a counter to reinitialize it
// the data is to refill the counters from disk
type Manager struct {
	db   *bolt.DB
	conf *config.Config
}

// NewManager ...
func NewManager() *Manager {
	conf := config.GetConfig()
	dataPath := conf.DataDir
	err := os.MkdirAll(dataPath, 0777)
	utils.PanicOnError(err)
	infoPath := filepath.Join(config.GetConfig().InfoDir, "system.db")
	db, err := bolt.Open(infoPath, 0777, nil)
	utils.PanicOnError(err)
	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("info")); err != nil {
			return err
		}
		_, err := tx.CreateBucketIfNotExists([]byte("domain"))
		return err
	})
	utils.PanicOnError(err)
	return &Manager{db, conf}
}

// GetFile ...
func (m *Manager) GetFile(id string) (*os.File, error) {
	f, err := os.OpenFile(filepath.Join(m.conf.DataDir, id), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (m *Manager) saveToBoltDB(tx *bolt.Tx, bucketID string, info map[string]interface{}) error {
	b := tx.Bucket([]byte(bucketID))

	// Remove deleted keys
	err := b.ForEach(func(k, v []byte) error {
		if _, ok := info[string(k)]; !ok {
			if err := b.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		// TODO: print something here
	}

	for k, v := range info {
		rawInfo, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("saving info: %s", err)
		}
		if err := b.Put([]byte(k), rawInfo); err != nil {
			return err
		}
	}
	return nil
}

// SaveInfo ...
func (m *Manager) SaveInfo(info map[string]*datamodel.Info) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		tmpInfo := make(map[string]interface{})
		for k, v := range info {
			tmpInfo[k] = v
		}
		return m.saveToBoltDB(tx, "info", tmpInfo)
	})
}

// SaveDomains ...
func (m *Manager) SaveDomains(info map[string][]string) error {
	return m.db.Update(func(tx *bolt.Tx) error {
		tmpInfo := make(map[string]interface{})
		for k, v := range info {
			tmpInfo[k] = v
		}
		return m.saveToBoltDB(tx, "domain", tmpInfo)
	})
}

// LoadAllInfo ...
func (m *Manager) LoadAllInfo() (map[string]*datamodel.Info, error) {
	infos := map[string]*datamodel.Info{}
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("info"))
		err := b.ForEach(func(k, v []byte) error {
			var info *datamodel.Info
			if err := json.Unmarshal(v, &info); err != nil {
				return err
			}
			infos[string(k)] = info
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// LoadAllDomains ...
func (m *Manager) LoadAllDomains() (map[string][]string, error) {
	ids := map[string][]string{}
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("domain"))
		err := b.ForEach(func(k, v []byte) error {
			var vals []string
			if err := json.Unmarshal(v, &vals); err != nil {
				return err
			}
			ids[string(k)] = vals
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// Close ...
func (m *Manager) Close() error {
	return m.db.Close()
}
