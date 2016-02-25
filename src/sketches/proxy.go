package sketches

import (
	"fmt"
	"sync"

	"github.com/njpatel/loggo"

	"datamodel"
)

var logger = loggo.GetLogger("sketches")

// SketchProxy ...
type SketchProxy struct {
	*datamodel.Info
	sketch datamodel.Sketcher
	lock   sync.RWMutex
}

// Add ...
func (sp *SketchProxy) Add(values [][]byte) (bool, error) {
	sp.lock.Lock()
	defer sp.lock.Unlock()
	return sp.sketch.Add(values)
}

// Get ...
func (sp *SketchProxy) Get(data interface{}) (interface{}, error) {
	switch datamodel.GetTypeString(sp.GetType()) {
	case datamodel.HLLPP:
		return sp.sketch.Get(nil)
	case datamodel.CML:
		return sp.sketch.Get(data)
	case datamodel.TopK:
		return sp.sketch.Get(nil)
	case datamodel.Bloom:
		return sp.sketch.Get(data)
	default:
		return nil, fmt.Errorf("Invalid sketch type: %s", sp.GetType())
	}
}

// CreateSketch ...
func CreateSketch(info *datamodel.Info) (*SketchProxy, error) {
	var err error
	var sketch datamodel.Sketcher
	sp := &SketchProxy{info, sketch, sync.RWMutex{}}

	switch datamodel.GetTypeString(info.GetType()) {
	case datamodel.HLLPP:
		sp.sketch, err = NewHLLPPSketch(info)
	case datamodel.CML:
		sp.sketch, err = NewCMLSketch(info)
	case datamodel.TopK:
		sp.sketch, err = NewTopKSketch(info)
	case datamodel.Bloom:
		sp.sketch, err = NewBloomSketch(info)
	default:
		return nil, fmt.Errorf("Invalid sketch type: %s", sp.GetType())
	}

	if err != nil {
		return nil, err
	}
	return sp, nil
}
