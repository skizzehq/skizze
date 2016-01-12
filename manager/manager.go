package manager

import (
	"fmt"
	"sort"
	"sync"
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
	domains  *domainManager
	lock     sync.RWMutex
	ticker   *time.Ticker
	storage  *storage.Manager
}

func (m *Manager) saveSketch(info *datamodel.Info) error {
	return m.sketches.save(info)
}

func (m *Manager) lockSketches() {
	for _, v := range m.infos.info {
		_ = m.sketches.lock(v)
	}
}

// Save ...
func (m *Manager) Save() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// 1) deep copy all sketches Info
	infoCopy := m.infos.getCopy()
	// 2) save DEFAULT SETTINGS
	// TODO: save defaut settings
	// 3) lock all sketches from being allowed to do ADD
	m.lockSketches()
	// 4) Clear AOF
	// TODO: clear AOF

	// 5) For each sketch
	var wg sync.WaitGroup
	running := 0
	for _, v := range m.infos.info {
		wg.Add(1)
		running++
		go func(info *datamodel.Info) {
			// a) save sketch
			if err := m.saveSketch(info); err != nil {
				// TODO: log something here
			}
			// b) replay from AOF (SELECT * FROM ops WHERE sketchId = ?)
			// TODO: Replay from AOF
			// c) unlock sketch
			v.Unlock()
			wg.Done()
		}(v)
		// Just 4 at a time
		if running%4 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()

	// 6) Save deep copied sketches info from previously
	m.infos.save(infoCopy)
	return nil
}

// NewManager ...
func NewManager() *Manager {
	storage := storage.NewManager()
	sketches := newSketchManager(storage)
	infos := newInfoManager(storage)
	domains := newDomainManager(infos, sketches, storage)

	m := &Manager{
		sketches: sketches,
		infos:    infos,
		domains:  domains,
		lock:     sync.RWMutex{},
		ticker:   time.NewTicker(time.Second * time.Duration(config.GetConfig().SaveThresholdSeconds)),
		storage:  storage,
	}

	for _, info := range infos.info {
		utils.PanicOnError(sketches.load(info))
	}

	// Set up saving on intervals
	go func() {
		for _ = range m.ticker.C {
			if m.Save() != nil {
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

// CreateDomain ...
func (m *Manager) CreateDomain(info *datamodel.Info) error {
	types := datamodel.GetTypes()
	infos := make(map[string]*datamodel.Info)
	for _, typ := range types {
		tmpInfo := datamodel.Info(*info)
		tmpInfo.Type = typ
		infos[tmpInfo.ID()] = &tmpInfo
	}
	return m.domains.create(info.Name, infos)
}

// AddToSketch ...
func (m *Manager) AddToSketch(info *datamodel.Info, values []string) error {
	return m.sketches.add(info.ID(), values)
}

// DeleteSketch ...
func (m *Manager) DeleteSketch(info *datamodel.Info) error {
	if err := m.infos.delete(info); err != nil {
		return err
	}
	return m.sketches.delete(info)
}

type getSketchesResults [][2]string

func (slice getSketchesResults) Len() int {
	return len(slice)
}

func (slice getSketchesResults) Less(i, j int) bool {
	if slice[i][0] == slice[j][0] {
		return slice[i][1] < slice[j][1]
	}
	return slice[i][0] < slice[j][0]
}

func (slice getSketchesResults) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// GetSketches return a list of sketch tuples [name, type]
func (m *Manager) GetSketches() [][2]string {
	sketches := getSketchesResults{}
	for _, v := range m.infos.info {
		sketches = append(sketches, [2]string{v.Name, v.Type})
	}
	sort.Sort(sketches)
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
