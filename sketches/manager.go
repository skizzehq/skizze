package sketches

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

// Manager is responsible for manipulating the sketches and syncing to disk
type Manager struct {
	sketches map[string]*SketchProxy
	info     map[string]*datamodel.Info
	lock     sync.RWMutex
	saving   uint64
	ticker   *time.Ticker
	storage  *storage.Manager
}

// CreateSketch ...
func (m *Manager) CreateSketch(info *datamodel.Info) error {
	if _, ok := m.info[info.ID()]; ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" already exists`,
			info.Type, info.Name)
	}
	sketch, err := createSketch(info)
	if err != nil {
		return err
	}
	m.sketches[info.ID()] = sketch
	m.info[info.ID()] = info
	return nil
}

// DeleteSketch ...
func (m *Manager) DeleteSketch(info *datamodel.Info) error {
	if _, ok := m.info[info.ID()]; !ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" does not exists`,
			info.Type, info.Name)
	}
	delete(m.sketches, info.ID())
	delete(m.info, info.ID())
	return nil
}

// GetSketches return a list of sketch tuples [name, type]
func (m *Manager) GetSketches() [][2]string {
	sketches := [][2]string{}
	for _, v := range m.info {
		sketches = append(sketches, [2]string{v.Name, v.Type})
	}
	return sketches
}

// AddToSketch ...
func (m *Manager) AddToSketch(info *datamodel.Info, values []string) error {
	count := m.saving
	if count > 0 {
		// FIXME: Add an AOF
		fmt.Println("can't add", count)
	}
	byts := make([][]byte, len(values), len(values))
	for i, v := range values {
		byts[i] = []byte(v)
	}
	// FIXME: return if adding was successful or not
	_, err := m.sketches[info.ID()].Add(byts)
	return err
}

// GetFromSketch ...
func (m *Manager) GetFromSketch(info *datamodel.Info, data interface{}) (interface{}, error) {
	var values []string
	if data != nil {
		values = data.([]string)
	}

	byts := make([][]byte, len(values), len(values))
	for i, v := range values {
		byts[i] = []byte(v)
	}
	v, ok := m.sketches[info.ID()]
	if !ok {
		return nil, fmt.Errorf("No such key %s", info.ID())
	}
	return v.Get(byts)
}

// Save ...
func (m *Manager) Save(infos map[string]*datamodel.Info) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(infos) == 0 {
		infos = m.info
	}

	// FIXME: Use chan instead
	var wg sync.WaitGroup
	running := 0
	for k, v := range infos {
		wg.Add(1)
		running++
		go func(ID string, info *datamodel.Info) {
			atomic.AddUint64(&m.saving, 1)
			if err := m.saveSketch(info); err != nil {
				// TODO: log something here
			}
			time.Sleep(time.Second * 2)
			wg.Done()
		}(k, v)
		// Just 4 at a time
		if running%4 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	m.saving = 0
	return nil
}

func (m *Manager) saveSketch(info *datamodel.Info) error {
	sketch, ok := m.sketches[info.ID()]
	if !ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" does not exists`,
			info.Type, info.Name)
	}

	file, err := m.storage.GetFile(info.ID())
	defer utils.CloseFile(file)

	if err != nil {
		return fmt.Errorf(`Could not get destination file for sketch of type "%s" and name "%s", %q`,
			info.Type, info.Name, err)
	}
	if err := sketch.Save(file); err != nil {
		return fmt.Errorf(`Could not save sketch "%s" with name "%s", %q`,
			info.Type, info.Name, err)
	}
	return m.storage.SaveInfo(map[string]*datamodel.Info{info.ID(): info})
}

func (m *Manager) loadAll() error {
	var err error
	if m.info, err = m.storage.LoadAllInfo(); err != nil {
		return err
	}

	for _, v := range m.info {
		if err := m.loadSketch(v); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) loadSketch(info *datamodel.Info) error {
	sketch, ok := m.sketches[info.ID()]
	if ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" already loaded`,
			info.Type, info.Name)
	}

	file, err := m.storage.GetFile(info.ID())
	defer utils.CloseFile(file)

	if err != nil {
		return fmt.Errorf(`Could not get find file for sketch of type "%s" and name "%s", %q`,
			info.Type, info.Name, err)
	}
	sketch, err = loadSketch(info, file)
	if err != nil {
		return fmt.Errorf(`Could not load sketch "%s" with name "%s", %q`,
			info.Type, info.Name, err)
	}
	m.sketches[info.ID()] = sketch
	return nil
}

func (m *Manager) autoSave() {
	for _ = range m.ticker.C {
		if m.Save(map[string]*datamodel.Info{}) != nil {
			// FIXME: print out something
		}
	}
}

// NewManager ...
func NewManager() (*Manager, error) {
	m := &Manager{
		make(map[string]*SketchProxy),
		make(map[string]*datamodel.Info),
		sync.RWMutex{},
		0,
		time.NewTicker(time.Second * time.Duration(config.GetConfig().SaveThresholdSeconds)),
		storage.NewManager(),
	}
	if err := m.loadAll(); err != nil {
		return nil, err
	}
	// Set up saving on intervals
	go m.autoSave()
	return m, nil
}

// Destroy ...
func (m *Manager) Destroy() {
	m.ticker.Stop()
	_ = m.storage.Close()
}
