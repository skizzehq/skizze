package sketches

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/seiflotfy/skizze/datamodel"
	"github.com/seiflotfy/skizze/utils"
)

// SketchProxy ...
type SketchProxy struct {
	*datamodel.Info
	sketch datamodel.Sketch
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
	switch sp.Type {
	case datamodel.HLLPP:
		return sp.sketch.Get(nil)
	case datamodel.CML:
		return sp.sketch.Get(data)
	case datamodel.TopK:
		return sp.sketch.Get(nil)
	case datamodel.Bloom:
		return sp.sketch.Get(data)
	default:
		return nil, errors.New("Invalid sketch type: " + sp.Type)
	}
}

// Save ...
func (sp *SketchProxy) Save(file *os.File) error {
	data, err := sp.sketch.Marshal()
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

// CreateSketch ...
func CreateSketch(info *datamodel.Info) (*SketchProxy, error) {
	var err error
	var sketch datamodel.Sketch
	sp := &SketchProxy{info, sketch, sync.RWMutex{}}

	switch info.Type {
	case datamodel.HLLPP:
		sp.sketch, err = NewHLLPPSketch(info)
	case datamodel.CML:
		sp.sketch, err = NewCMLSketch(info)
	case datamodel.TopK:
		sp.sketch, err = NewTopKSketch(info)
	case datamodel.Bloom:
		sp.sketch, err = NewBloomSketch(info)
	default:
		return nil, errors.New("Invalid sketch type: " + info.Type)
	}

	if err != nil {
		return nil, err
	}
	return sp, nil
}

// LoadSketch ...
func LoadSketch(info *datamodel.Info, file *os.File) (*SketchProxy, error) {
	var sketch datamodel.Sketch
	sp := &SketchProxy{info, sketch, sync.RWMutex{}}

	size, err := utils.GetFileSize(file)
	if err != nil {
		return nil, err
	}

	data := make([]byte, size, size)
	_, err = file.Read(data)
	if err != nil {
		return nil, fmt.Errorf("Error loading data for sketch: %s", info.ID())
	}

	switch info.Type {
	case datamodel.HLLPP:
		sp.sketch = &HLLPPSketch{}
		err = sp.sketch.Unmarshal(info, data)
	case datamodel.CML:
		sp.sketch = &CMLSketch{}
		err = sp.sketch.Unmarshal(info, data)
	case datamodel.TopK:
		sp.sketch = &TopKSketch{}
		err = sp.sketch.Unmarshal(info, data)
	case datamodel.Bloom:
		sp.sketch = &BloomSketch{}
		err = sp.sketch.Unmarshal(info, data)
	default:
		//logger.Info.Println("Invalid sketch type", info.Type)
	}
	if err != nil {
		return nil, err
	}
	return sp, nil
}
