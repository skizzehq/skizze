package manager

import (
	"fmt"

	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/sketches"
	"github.com/seiflotfy/skizze/storage"
	"github.com/seiflotfy/skizze/utils"
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
		return fmt.Errorf(`Could not load sketch "%s" with name "%s", %q`,
			info.Type, info.Name, err)
	}
	m.sketches[info.ID()] = sketch
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

func (m *sketchManager) save(info *datamodel.Info) error {
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

func (m *sketchManager) add(id string, values []string) error {
	byts := make([][]byte, len(values), len(values))
	for i, v := range values {
		byts[i] = []byte(v)
	}
	// FIXME: return if adding was successful or not
	_, err := m.sketches[id].Add(byts)
	return err
}

func (m *sketchManager) delete(info *datamodel.Info) error {
	if _, ok := m.sketches[info.ID()]; !ok {
		return fmt.Errorf(`Sketch of type "%s" with name "%s" does not exists`,
			info.Type, info.Name)
	}
	delete(m.sketches, info.ID())
	return nil
}

func (m *sketchManager) get(info *datamodel.Info, data interface{}) (interface{}, error) {
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
