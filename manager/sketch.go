package manager

import (
	"fmt"

	"github.com/skizzehq/skizze/datamodel"
	"github.com/skizzehq/skizze/sketches"
	"github.com/skizzehq/skizze/storage"
	"github.com/skizzehq/skizze/utils"
)

type sketchManager struct {
	storage  *storage.Manager
	sketches map[string]*sketches.SketchProxy
}

func newSketchManager(storage *storage.Manager) *sketchManager {
	return &sketchManager{
		sketches: make(map[string]*sketches.SketchProxy),
		storage:  storage,
	}
}

func (m *sketchManager) load(info *datamodel.Info) error {
	sketch, ok := m.sketches[info.ID()]
	if ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" already loaded`,
			info.Type, info.Name)
	}

	file, err := m.storage.GetFile(info.ID())
	utils.PanicOnError(err)
	defer utils.CloseFile(file)

	if err != nil {
		return fmt.Errorf(`Could not get find file for sketch of type "%s" and name "%s", %q`,
			info.Type, info.Name, err)
	}
	sketch, err = sketches.LoadSketch(info, file)
	if err != nil {
		return fmt.Errorf(`Could not load sketch "%s" with name "%s", %q`, info.Type, info.Name, err)
	}
	m.sketches[info.ID()] = sketch
	return nil
}

func (m *sketchManager) setLockAll(b bool) {
	for _, sketch := range m.sketches {
		if b {
			sketch.Lock()
		} else {
			sketch.Unlock()
		}
	}
}

func (m *sketchManager) setLock(id string, b bool) error {
	sketch, ok := m.sketches[id]
	if !ok {
		return fmt.Errorf(`Sketch "%s" does not exist`, id)
	}
	if b {
		sketch.Lock()
	} else {
		sketch.Unlock()
	}
	return nil
}

// CreateSketch ...
func (m *sketchManager) create(info *datamodel.Info) error {
	sketch, err := sketches.CreateSketch(info)
	if err != nil {
		return err
	}
	m.sketches[info.ID()] = sketch
	return nil
}

func (m *sketchManager) save(id string) error {
	sketch, ok := m.sketches[id]
	if !ok {
		return fmt.Errorf(`Sketch "%s" does not exists`, id)
	}

	file, err := m.storage.GetFile(id)
	defer utils.CloseFile(file)

	if err != nil {
		return fmt.Errorf(`Could not get destination file for sketch "%s", %q`, id, err)
	}
	if err := sketch.Save(file); err != nil {
		return fmt.Errorf(`Could not save sketch "%s", %q`, id, err)
	}
	return nil
}

func (m *sketchManager) add(id string, values []string) error {
	sketch, ok := m.sketches[id]
	if !ok {
		return fmt.Errorf(`Sketch "%s" does not exists`, id)
	}
	if sketch.Locked() {
		// Append to File here
		return nil //&lockedError{}
	}

	byts := make([][]byte, len(values), len(values))
	for i, v := range values {
		byts[i] = []byte(v)
	}
	// FIXME: return if adding was successful or not
	_, err := sketch.Add(byts)
	return err
}

func (m *sketchManager) delete(id string) error {
	if _, ok := m.sketches[id]; !ok {
		return fmt.Errorf(`Sketch "%s" does not exists`, id)
	}
	delete(m.sketches, id)
	return nil
}

func (m *sketchManager) get(id string, data interface{}) (interface{}, error) {
	var values []string
	if data != nil {
		values = data.([]string)
	}
	byts := make([][]byte, len(values), len(values))
	for i, v := range values {
		byts[i] = []byte(v)
	}
	v, ok := m.sketches[id]
	if !ok {
		return nil, fmt.Errorf("No such key %s", id)
	}
	return v.Get(byts)
}
