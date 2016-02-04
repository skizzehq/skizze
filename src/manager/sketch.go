package manager

import (
	"fmt"

	"datamodel"
	"sketches"
)

type sketchManager struct {
	sketches map[string]*sketches.SketchProxy
}

func newSketchManager() *sketchManager {
	return &sketchManager{
		sketches: make(map[string]*sketches.SketchProxy),
	}
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
