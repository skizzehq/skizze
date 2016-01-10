package manager

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/seiflotfy/skizze/config"
	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
)

func isValidType(info *datamodel.Info) bool {
	return info.Type == datamodel.Bloom || info.Type == datamodel.CML ||
		info.Type == datamodel.HLLPP || info.Type == datamodel.TopK
}

// Manager is responsible for manipulating the sketches and syncing to disk
type Manager struct {
	infos    *infoManager
	sketches *sketchManager
	lock     sync.RWMutex
	saving   uint64
	ticker   *time.Ticker
	storage  *storage.Manager
}

func (m *Manager) saveSketch(info *datamodel.Info) error {
	atomic.AddUint64(&m.saving, 1)
	if err := m.sketches.save(info); err != nil {
		return err
	}
	return m.infos.save(info.ID())
}

// Save ...
func (m *Manager) Save(infos map[string]*datamodel.Info) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(infos) == 0 {
		infos = m.infos.info
	}

	// FIXME: Use chan instead
	var wg sync.WaitGroup
	running := 0
	for _, v := range infos {
		wg.Add(1)
		running++
		go func(info *datamodel.Info) {
			if err := m.saveSketch(info); err != nil {
				// TODO: log something here
			}
			wg.Done()
		}(v)
		// Just 4 at a time
		if running%4 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	m.saving = 0
	return nil
}

// NewManager ...
func NewManager() *Manager {
	storage := storage.NewManager()
	sManager := newSketchManager(storage)
	iManager := newInfoManager(storage)

	m := &Manager{
		sketches: sManager,
		infos:    iManager,
		lock:     sync.RWMutex{},
		saving:   0,
		ticker:   time.NewTicker(time.Second * time.Duration(config.GetConfig().SaveThresholdSeconds)),
		storage:  storage,
	}

	for _, info := range iManager.info {
		utils.PanicOnError(sManager.load(info))
	}

	// Set up saving on intervals
	go func() {
		for _ = range m.ticker.C {
			if m.Save(map[string]*datamodel.Info{}) != nil {
				// FIXME: print out something
			}
		}
	}()
	return m
}

// CreateSketch ...
func (m *Manager) CreateSketch(info *datamodel.Info) error {
	if !isValidType(info) {
		return fmt.Errorf("Can not create sketch of type %s, invalid type.", info.Type)
	}
	if err := m.infos.create(info); err != nil {
		return err
	}
	if err := m.sketches.create(info); err != nil {
		// If error occurred during creation of sketch, delete info
		if err2 := m.infos.delete(info); err2 != nil {
			return fmt.Errorf("%q\n%q ", err, err2)
		}
		return err
	}
	return nil
}

// AddToSketch ...
func (m *Manager) AddToSketch(info *datamodel.Info, values []string) error {
	count := m.saving
	if count > 0 {
		// FIXME: Add an AOF
		fmt.Println("can't add", count)
	}
	return m.sketches.add(info.ID(), values)
}

// DeleteSketch ...
func (m *Manager) DeleteSketch(info *datamodel.Info) error {
	if err := m.infos.delete(info); err != nil {
		return err
	}
	return m.sketches.delete(info)
}

// GetSketches return a list of sketch tuples [name, type]
func (m *Manager) GetSketches() [][2]string {
	sketches := [][2]string{}
	for _, v := range m.infos.info {
		sketches = append(sketches, [2]string{v.Name, v.Type})
	}
	return sketches
}

// GetFromSketch ...
func (m *Manager) GetFromSketch(info *datamodel.Info, data interface{}) (interface{}, error) {
	return m.sketches.get(info, data)
}

// Destroy ...
func (m *Manager) Destroy() {
	m.ticker.Stop()
	_ = m.storage.Close()
}
